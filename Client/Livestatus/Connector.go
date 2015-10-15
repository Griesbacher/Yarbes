package Livestatus

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
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

	length := 1
	for length > 0 {
		message, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				resultError <- err
				continue
			}
		}
		length = len(message)
		if length > 0 {
			csvReader := csv.NewReader(strings.NewReader(string(message)))
			csvReader.Comma = ';'
			records, err := csvReader.Read()
			if err != nil {
				resultError <- err
				continue
			}
			result <- records
		}
	}
	close(result)
	close(resultError)
}
