package Livestatus

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"net"
)

//Connector fetches data from Livestatussocket
type Connector struct {
	LivestatusAddress string
	ConnectionType    string
}

//connectToLivestatus queries livestatus and sends the result line by line in the result channel
func (connector Connector) connectToLivestatus(query string, result chan []string, resultError chan error) {
	var conn net.Conn
	switch connector.ConnectionType {
	case "tcp":
		conn, _ = net.Dial("tcp", connector.LivestatusAddress)
	case "file":
		conn, _ = net.Dial("unix", connector.LivestatusAddress)
	default:
		resultError <- errors.New("Connection type is unkown, options are: tcp, file. Input:" + connector.ConnectionType)
		close(result)
		close(resultError)
		return
	}
	if conn == nil {
		resultError <- errors.New("Unable to connect to livestatus" + connector.LivestatusAddress)
		close(result)
		close(resultError)
		return
	}

	defer conn.Close()
	fmt.Fprintf(conn, query)
	reader := bufio.NewReader(conn)
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\x02'
	records, err := csvReader.ReadAll()
	if err != nil {
		resultError <- err
	}
	for _, record := range records{
		result <- record
	}
	close(result)
	close(resultError)
}
