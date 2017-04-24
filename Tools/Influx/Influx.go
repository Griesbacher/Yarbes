package Influx

import "github.com/influxdata/influxdb/client/v2"

//QueryDB need a influxdb client and a command which will be send to the inflxdb
func QueryDB(clnt client.Client, cmd, db string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}
