package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Logging"
	"github.com/griesbacher/Yarbes/NetworkInterfaces/Outgoing"
	"github.com/griesbacher/Yarbes/Tools/Strings"
	"github.com/influxdb/influxdb/client/v2"
	"github.com/influxdb/influxdb/models"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

//TablePrefix is the prefix for all new influxdb tables
const TablePrefix = "spam_"

var table string
var logger *Logging.Client
var jsonMap map[string]interface{}
var timestampField string
var resultField string

func main() {
	var moduleConfigPath string
	var clientConfigPath string
	var tableName string
	var delay int64
	var level int64
	var save bool
	var messageField string
	var event string
	flag.Usage = func() {
		fmt.Println(`Yarbes-EventPerTime by Philip Griesbacher @ 2015
Commandline Parameter:
-moduleConfigPath Path to the module config file. If no file path is given the default is ./Module/External/EventsPerTime/EventsPerTime.gcfg.
-clienteConfigPath Path to the client config file. If no file path is given the default is ./clientConfig.gcfg.`)
	}
	flag.StringVar(&moduleConfigPath, "moduleConfigPath", "Module/External/EventsPerTime/EventsPerTime.gcfg", "path to the module config file")
	flag.StringVar(&clientConfigPath, "clienteConfigPath", "clientConfig.gcfg", "path to the client config file")
	flag.StringVar(&tableName, "tableName", "t", "tablenamen which will be used to save the statistics")
	flag.Int64Var(&delay, "delay", 5, "delay in seconds")
	flag.Int64Var(&level, "level", 2, "when this level is reached a composit event will be created, events per delay")
	flag.BoolVar(&save, "save", false, "true for the original event false for the dummy delay event")
	flag.StringVar(&messageField, "messageField", "message", "references the filed in which the message can be found")
	flag.StringVar(&timestampField, "timestampField", "time", "references the filed in which the unixtimestamp can be found")
	flag.StringVar(&event, "event", "", "the event")
	flag.StringVar(&resultField, "returnField", "EventPerTimeResult", "the result, stored messages")
	flag.Parse()
	//Load Config
	Config.InitEventsPerTimeConfig(moduleConfigPath)
	Config.InitClientConfig(clientConfigPath)

	//create influxdb client
	u, err := url.Parse(Config.GetEventsPerTimeConfig().InfluxDB.Server)
	if err != nil {
		panic(err)
	}
	c := client.NewClient(client.Config{
		URL:      u,
		Username: Config.GetEventsPerTimeConfig().InfluxDB.Username,
		Password: Config.GetEventsPerTimeConfig().InfluxDB.Password,
	})

	rand.Seed(time.Now().Unix())

	table = TablePrefix + tableName

	_, err = queryDB(c, fmt.Sprintf("CREATE DATABASE %s", Config.GetEventsPerTimeConfig().InfluxDB.Database))
	if err != nil {
		panic(err)
	}

	logger, err = Logging.NewClient(Config.GetClientConfig().LogServer.RPCInterface)
	if err != nil {
		logger = Logging.NewLocalClient()
	}

	var jsonInterface interface{}
	err = json.Unmarshal([]byte(event), &jsonInterface)
	if err != nil {
		panic(err)
	}
	jsonMap = jsonInterface.(map[string]interface{})
	var eventTimestamp int
	switch value := jsonMap[timestampField].(type) {
	case string:
		eventTimestamp, err = strconv.Atoi(value)
	case float64:
		tmpStamp := jsonMap[timestampField].(float64)
		eventTimestamp = int(tmpStamp)
	default:
		err = errors.New("Type unkown")
	}

	if err != nil {
		logger.Error(err)
		os.Exit(10)
	}

	//TODO: durch filelock absichern
	if save {
		eventRPC := Outgoing.NewRPCInterface(Config.GetClientConfig().Backend.RPCInterface)
		err = eventRPC.Connect()
		if err != nil {
			logger.Error(err)
			os.Exit(11)
		}
		eventMessage := jsonMap[messageField].(string)
		addEvent(c, eventTimestamp, eventMessage)

		delay := time.Duration(delay) * time.Second
		eventRPC.CreateDelayedEvent([]byte(`{"type":"EventsPerTime","`+timestampField+`":`+fmt.Sprint(eventTimestamp)+`}`), &delay)
	} else {
		handlePoint(c, eventTimestamp, int(delay), int(level))
	}
}

var globalCount = 0

func handlePoint(c client.Client, eventTimestamp, timerange, level int) {
	min := eventTimestamp - timerange
	if min < 0 {
		min = 0
	}
	max := eventTimestamp + timerange

	q := fmt.Sprintf("Select handled, msg from %s WHERE time = %ds; Select handled, msg, count from %s WHERE time > %ds and time < %ds", table, eventTimestamp, table, min, max)
	res, err := queryDB(c, q)
	if err != nil {
		panic(err)
	} else {
		globalCount = len(res[1].Series[0].Values)
		if !res[0].Series[0].Values[0][1].(bool) {
			if globalCount < level {
				setOwnPointsHandled(c, res[1].Series[0], time.Unix(int64(eventTimestamp), int64(0)).UTC())
			} else {
				setPointRangeHandled(c, res[1].Series[0])
			}
		}
	}
}

func notify(a ...interface{}) {
	msg := `"`
	for _, message := range a {
		msg += strings.Trim(fmt.Sprint(message), `\"`)
	}
	msg += `"`
	jsonMap[resultField] = msg
	jsonMap[resultField+"count"] = globalCount
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}
	fmt.Println(`{"Event": ` + string(jsonBytes) + `}`)
}

func setOwnPointsHandled(c client.Client, row models.Row, eventTimestamp time.Time) {
	bp := genBatchPoints()
	countID := Strings.IndexOf(row.Columns, "count")
	msgID := Strings.IndexOf(row.Columns, "msg")
	handledID := Strings.IndexOf(row.Columns, "handled")
	for _, value := range row.Values {
		if value[handledID] == true {
			continue
		}
		timeStamp, err := time.Parse(time.RFC3339, fmt.Sprint(value[0]))
		if err != nil {
			panic(err)
		}
		if timeStamp == eventTimestamp {
			fields := map[string]interface{}{
				"handled": true,
				"msg":     value[msgID],
			}
			bp.AddPoint(client.NewPoint(table, map[string]string{"count": fmt.Sprint(value[countID])}, fields, timeStamp))
			notify(value[2])
		}
		break
	}
	writePoints(c, bp)
}

func setPointRangeHandled(c client.Client, row models.Row) {
	bp := genBatchPoints()
	countID := Strings.IndexOf(row.Columns, "count")
	msgID := Strings.IndexOf(row.Columns, "msg")
	handledID := Strings.IndexOf(row.Columns, "handled")
	message := []interface{}{}
	for _, value := range row.Values {
		if value[handledID] == true {
			continue
		}
		fields := map[string]interface{}{
			"handled": true,
			"msg":     value[msgID],
		}
		timeStamp, err := time.Parse(time.RFC3339, fmt.Sprint(value[0]))
		if err != nil {
			panic(err)
		}
		bp.AddPoint(client.NewPoint(table, map[string]string{"count": fmt.Sprint(value[countID])}, fields, timeStamp))
		message = append(message, value[msgID])

	}
	writePoints(c, bp)
	notify(message)
}

func genEvents(c client.Client) {
	queryDB(c, "DROP SERIES FROM "+table)
	for _, i := range genRange() {
		for j := 0; j < 2; j++ {
			addEvent(c, i, "Hallo")
		}
	}
}

func addEvent(c client.Client, timestamp int, msg string) {
	count := fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int63())
	bp := genBatchPoints()
	fields := map[string]interface{}{
		"msg":     msg,
		"handled": false,
	}
	bp.AddPoint(client.NewPoint(table, map[string]string{"count": count}, fields, time.Unix(int64(timestamp), int64(0))))
	writePoints(c, bp)

}

func genRange() []int {
	result := []int{}
	for j := 0.0; j < 32; j += 2 {
		diff := math.Sin(j/10) * 100
		length := len(result)
		if length == 0 {
			result = append(result, 1)
		} else {
			result = append(result, result[length-1]+int(math.Abs(diff-float64(result[length-1]))))
		}
	}
	return result
}

func writePoints(c client.Client, bp client.BatchPoints) {
	if c.Write(bp) != nil {
		panic(bp)
	}
}

func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: Config.GetEventsPerTimeConfig().InfluxDB.Database,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}

func genBatchPoints() client.BatchPoints {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{Database: Config.GetEventsPerTimeConfig().InfluxDB.Database, Precision: "us"})
	if err != nil {
		panic(err)
	}
	return bp
}
