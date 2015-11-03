package LogServer

import (
	"fmt"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/influxdb/influxdb/client/v2"
	"net/url"
)

//Server receives LogMessages and does something with them
type Server struct {
	LogQueue     chan LogMessage
	quit         chan bool
	isRunning    bool
	influxClient client.Client
}

//TableName is the name of the table within influxdb
const TableName = "logs"

//NewLogServer constructs a new LogServer
func NewLogServer() *Server {
	var influxClient client.Client
	if Config.GetServerConfig().LogServer.InfluxAddress != "" {
		u, err := url.Parse(Config.GetServerConfig().LogServer.InfluxAddress)
		if err != nil {
			panic(err)
		}
		influxClient = client.NewClient(client.Config{
			URL:      u,
			Username: Config.GetServerConfig().LogServer.InfluxUsername,
			Password: Config.GetServerConfig().LogServer.InfluxPassword,
		})
		_, err = queryDB(influxClient, fmt.Sprintf("CREATE DATABASE %s", Config.GetServerConfig().LogServer.InfluxDatabase))
		if err != nil {
			panic(err)
		}
	}
	return &Server{LogQueue: make(chan LogMessage, 100), quit: make(chan bool), isRunning: false, influxClient: influxClient}
}

//Start starts the LogServer
func (log Server) Start() {
	if !log.isRunning {
		go log.handleLog()
	}
}

//Stop stops the LogServer
func (log Server) Stop() {
	if log.isRunning {
		log.quit <- true
		<-log.quit
	}
}

//IsRunning returns true if the daemon is running
func (log Server) IsRunning() bool {
	return log.isRunning
}

func (log *Server) handleLog() {
	log.isRunning = true
	var message LogMessage
	for {
		select {
		case <-log.quit:
			log.quit <- true
			return
		case message = <-log.LogQueue:
			fmt.Println(message.String())
			go log.logToInfluxDB(message)
		}
	}
}

func (log Server) logToInfluxDB(message LogMessage) {
	if &log.influxClient == nil {
		return
	}
	bp, err := genBatchPoints()
	if err != nil {
		fmt.Println(err)
	}
	fields := map[string]interface{}{
		"msg":        message.Message,
		"event":      message.Event.String(),
		"source":     message.Source,
		"serveritry": SeverityToString(message.Severity),
	}
	point, err := client.NewPoint(TableName, map[string]string{}, fields, message.Timestamp)
	if err != nil {
		fmt.Println(err)
	}
	bp.AddPoint(point)

	err = log.influxClient.Write(bp)
	if err != nil {
		fmt.Println(err)
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

func genBatchPoints() (client.BatchPoints, error) {
	return client.NewBatchPoints(client.BatchPointsConfig{Database: Config.GetServerConfig().LogServer.InfluxDatabase, Precision: "us"})
}
