package Livestatus

import (
	"fmt"
	"github.com/griesbacher/SystemX/Client"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/Logging"
	"github.com/griesbacher/SystemX/Tools/Strings"
	"time"
)

//Collector data from Livestatus
type Collector struct {
	conn      Connector
	quit      chan bool
	isRunning bool
	logger    Logging.Client
	creator   Client.EventCreatable
	converter *livestatusResultConverter
}

//LogQuery will be used for every livestatus query
//TODO: fieldseperator durch 1 2 5 6 ersetzen
const LogQuery = `GET log
Columns: type time host_name current_service_display_name long_plugin_output
Filter: time > %d
WaitTrigger: log
WaitTimeout: 10000
OutputFormat: csv
Separators: 10 2 5 6

`

//NewCollector constructs a new Livestatus Collector
func NewCollector(logger Logging.Client, eventCreator Client.EventCreatable) *Collector {
	connector := Connector{LivestatusAddress: Config.GetClientConfig().Livestatus.Address, ConnectionType: Config.GetClientConfig().Livestatus.Type}
	converter := newLivestatusResultConverter(LogQuery)
	return &Collector{conn: connector, quit: make(chan bool), isRunning: false, logger: logger, creator: eventCreator, converter: converter}
}

//Start starts the collector
func (collector Collector) Start() {
	if !collector.isRunning {
		go collector.work()
	}
}

//Stop stops the collector
func (collector Collector) Stop() {
	collector.quit <- true
	<-collector.quit
}

func (collector *Collector) work() {
	collector.isRunning = true
	start := time.Now()
	var result chan []string
	var errorChan chan error
	var oldCache []string
	var newCache []string
	oldCache = []string{}
	for {
		select {
		case <-collector.quit:
			collector.quit <- true
			return
		default:
			result = make(chan []string, 10)
			errorChan = make(chan error)
			newCache = []string{}
			timeToHandleRequest := time.Now().Sub(start)
			//fmt.Println(time.Now().Unix(), timeToHandleRequest)
			go collector.conn.connectToLivestatus(addTimestampToLivestatusQuery(LogQuery, timeToHandleRequest), result, errorChan)
			start = time.Now()
			queryRunning := true
			for queryRunning {
				select {
				case line, alive := <-result:
					if alive {
						newCache = append(newCache, fmt.Sprint(line))
						if !Strings.Contains(oldCache, fmt.Sprint(line)) {
							fmt.Println("-->", line)
							jsonEvent := collector.convertQueryResultToJSON(line)
							collector.sendEvent(jsonEvent)
						}
					} else {
						queryRunning = false
					}
				case err, alive := <-errorChan:
					if alive {
						fmt.Println(err, alive)
						collector.logger.Error(err)
					} else {
						queryRunning = false
					}
				case <-time.After(time.Duration(15) * time.Second):
					collector.logger.Debug("Livestatus collector timed out")
				case <-collector.quit:
					collector.quit <- true
					return
				}
			}
			oldCache = newCache
		}
	}
}

func (collector Collector) sendEvent(event []byte) {
	if len(event) == 0 {
		return
	}
	err := collector.creator.CreateEvent(event)
	if err != nil {
		collector.logger.Error(err)
	}
}

func addTimestampToLivestatusQuery(query string, durration time.Duration) string {
	return fmt.Sprintf(query, time.Now().Add((durration+time.Duration(1)*time.Second)*-2).Unix())
}

func (collector Collector) convertQueryResultToJSON(queryLine []string) []byte {
	event := collector.converter.createObject(queryLine)
	newEvent, err := Event.NewEventFromInterface(event)
	if err != nil {
		collector.logger.Error(err)
	}
	return newEvent.GetDataAsBytes()
}
