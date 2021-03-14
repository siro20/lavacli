// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"encoding/base64"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

type JobsListing struct {
	Description string `xmlrpm:"description"`
	DeviceType  string `xmlrpc:"device_type"`
	Health      string `xmlrpc:"health"`
	ID          int    `xmlrpc:"id"`
	State       string `xmlrpc:"state"`
	Submitter   string `xmlrpc:"submitter"`
}

func (c Connection) JobsList(state string, health string, start int, limit int) ([]JobsListing, error) {
	var ret []JobsListing

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

type JobState struct {
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

func (c Connection) JobsShow(id int) (*JobState, error) {
	var ret JobState

	err := c.con.Call("scheduler.jobs.show", id, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type JobDefintion string

func (c Connection) JobsDefinition(id int) (JobDefintion, error) {
	var ret JobDefintion

	err := c.con.Call("scheduler.jobs.definition", id, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

type JobErrors map[string]interface{}

func (c Connection) JobsValidate(def string, strict bool) (JobErrors, error) {
	var ret JobErrors
	var args []interface{}
	args = append(args, def)
	args = append(args, strict)

	err := c.con.Call("scheduler.jobs.validate", args, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c Connection) JobsSubmit(def string) ([]int, error) {
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

func (c Connection) JobsCancel(id int) error {

	err := c.con.Call("scheduler.jobs.cancel", id, nil)

	return err
}

func (c Connection) JobsFail(id int) error {

	err := c.con.Call("scheduler.jobs.fail", id, nil)

	return err
}

type JobLogsDecoded struct {
	DateTime string      `yaml:"dt"`
	Level    string      `yaml:"lvl"`
	Message  interface{} `yaml:"msg"`
}

type JobsLogs struct {
	Finished bool
	Data     string
	Decoded  []JobLogsDecoded
}

func (c Connection) JobsLogs(id int, raw bool) (*JobsLogs, error) {
	var ret []interface{}
	var ret2 JobsLogs

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
