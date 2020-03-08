// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"time"
)

type LavaJobsListing struct {
	Description string `xmlrpm:"description"`
	DeviceType  string `xmlrpc:"device_type"`
	Health      string `xmlrpc:"health"`
	ID          int    `xmlrpc:"id"`
	State       string `xmlrpc:"state"`
	Submitter   string `xmlrpc:"submitter"`
}

func (c LavaConnection) LavaJobsList() ([]LavaJobsListing, error) {
	var ret []LavaJobsListing

	err := c.con.Call("scheduler.jobs.list", nil, &ret)
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

type LavaJobIDs []int

func (c LavaConnection) LavaJobsSubmit(def string) (LavaJobIDs, error) {
	var ret LavaJobIDs

	err := c.con.Call("scheduler.jobs.submit", def, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c LavaConnection) LavaJobsCancel(id int) error {

	err := c.con.Call("scheduler.jobs.definition", id, nil)

	return err
}
