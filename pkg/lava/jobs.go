// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"time"

	"github.com/kolo/xmlrpc"
)

type LavaJobsListing struct {
	Description string `xmlrpm:"description"`
	DeviceType  string `xmlrpc:"device_type"`
	Health      string `xmlrpc:"health"`
	ID          int    `xmlrpc:"id"`
	State       string `xmlrpc:"state"`
	Submitter   string `xmlrpc:"submitter"`
}

func LavaJobsList(con *xmlrpc.Client) ([]LavaJobsListing, error) {
	var ret []LavaJobsListing

	err := con.Call("scheduler.jobs.list", nil, &ret)
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

func LavaJobsShow(con *xmlrpc.Client, id int) (*LavaJobState, error) {
	var ret LavaJobState

	err := con.Call("scheduler.jobs.show", id, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type LavaJobDefintion string

func LavaJobsDefinition(con *xmlrpc.Client, id int) (LavaJobDefintion, error) {
	var ret LavaJobDefintion

	err := con.Call("scheduler.jobs.definition", id, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

type LavaJobErrors map[string]interface{}

func LavaJobsValidate(con *xmlrpc.Client, def string, strict bool) (LavaJobErrors, error) {
	var ret LavaJobErrors
	var args []interface{}
	args = append(args, def)
	args = append(args, strict)

	err := con.Call("scheduler.jobs.validate", args, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type LavaJobIDs []int

func LavaJobsSubmit(con *xmlrpc.Client, def string) (LavaJobIDs, error) {
	var ret LavaJobIDs

	err := con.Call("scheduler.jobs.submit", def, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func LavaJobsCancel(con *xmlrpc.Client, id int) error {

	err := con.Call("scheduler.jobs.definition", id, nil)

	return err
}
