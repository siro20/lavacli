// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"encoding/base64"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

type LavaJobsListing struct {
	Description string `xmlrpm:"description"`
	DeviceType  string `xmlrpc:"device_type"`
	Health      string `xmlrpc:"health"`
	ID          int    `xmlrpc:"id"`
	State       string `xmlrpc:"state"`
	Submitter   string `xmlrpc:"submitter"`
}

func (c LavaConnection) LavaJobsList(state string, health string, start int, limit int) ([]LavaJobsListing, error) {
	var ret []LavaJobsListing

	var args []interface{}
	args = append(args, state)
	args = append(args, health)
	args = append(args, start)
	args = append(args, limit)

	err := c.con.Call("scheduler.jobs.list", args, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type LavaJobState struct {
	Description    string    `xmlrpc:"description"`
	DeviceType     string    `xmlrpc:"device_type"`
	Device         string    `xmlrpc:"device"`
	State          string    `xmlrpc:"state"`
	ID             int       `xmlrpc:"id"`
	EndTime        time.Time `xmlrpc:"end_time"`
	SubmitTime     time.Time `xmlrpc:"submit_time"`
	FailureComment string    `xmlrpc:"failure_comment"`
	Status         int       `xmlrpc:"status"`
	HealthCheck    bool      `xmlrpc:"health_check"`
	Pipeline       bool      `xmlrpc:"pipeline"`
	Tags           []string  `xmlrpc:"tags"`
	Visibility     string    `xmlrpc:"visibility"`
	Submitter      string    `xmlrpc:"submitter"`
	StartTime      time.Time `xmlrpc:"start_time"`
	Health         string    `xmlrpc:"health"`
}

func (c LavaConnection) LavaJobsShow(id int) (*LavaJobState, error) {
	var ret LavaJobState

	err := c.con.Call("scheduler.jobs.show", id, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type LavaJobDefintion string

func (c LavaConnection) LavaJobsDefinition(id int) (LavaJobDefintion, error) {
	var ret LavaJobDefintion

	err := c.con.Call("scheduler.jobs.definition", id, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

type LavaJobErrors map[string]interface{}

func (c LavaConnection) LavaJobsValidate(def string, strict bool) (LavaJobErrors, error) {
	var ret LavaJobErrors
	var args []interface{}
	args = append(args, def)
	args = append(args, strict)

	err := c.con.Call("scheduler.jobs.validate", args, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c LavaConnection) LavaJobsSubmit(def string) ([]int, error) {
	var ret []int
	var xmlRet interface{}

	err := c.con.Call("scheduler.jobs.submit", def, &xmlRet)
	if err != nil {
		return nil, err
	}

	switch x := xmlRet.(type) {
	case []int:
		for _, i := range x {
			ret = append(ret, int(x[i]))
		}
	case []int64:
		for _, i := range x {
			ret = append(ret, int(x[i]))
		}
	case int:
		ret = append(ret, int(x))
	case int64:
		ret = append(ret, int(x))
	default:
		return nil, fmt.Errorf("Got unexpected type: %T", x)
	}

	return ret, nil
}

func (c LavaConnection) LavaJobsCancel(id int) error {

	err := c.con.Call("scheduler.jobs.definition", id, nil)

	return err
}

type LavaJobLogsDecoded struct {
	DateTime string      `yaml:"dt"`
	Level    string      `yaml:"lvl"`
	Message  interface{} `yaml:"msg"`
}

type LavaJobsLogs struct {
	Finished bool
	Data     string
	Decoded  []LavaJobLogsDecoded
}

func (c LavaConnection) LavaJobsLogs(id int, raw bool) (*LavaJobsLogs, error) {
	var ret []interface{}
	var ret2 LavaJobsLogs

	err := c.con.Call("scheduler.jobs.logs", id, &ret)
	if err != nil {
		return nil, err
	}
	if len(ret) != 2 {
		return nil, fmt.Errorf("Unexpected server response")
	}
	finished, ok := ret[0].(bool)
	if !ok {
		return nil, fmt.Errorf("Invalid server response")
	}
	ret2.Finished = finished

	data, ok := ret[1].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid server response")
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	ret2.Data = string(decoded)

	if raw {
		return &ret2, nil
	}

	err = yaml.Unmarshal(decoded, &ret2.Decoded)
	if err != nil {
		return nil, err
	}
	return &ret2, nil
}
