// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"encoding/base64"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

// JobsListing represents data as returned by LAVA XMLRPC scheduler.jobs.list
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

// JobState represents data as returned by LAVA XMLRPC scheduler.jobs.show
type JobState struct {
	Description    string    `xmlrpc:"description" yaml:"description" json:"description"`
	DeviceType     string    `xmlrpc:"device_type"  yaml:"device_type" json:"device_type"`
	Device         string    `xmlrpc:"device"  yaml:"device" json:"device"`
	State          string    `xmlrpc:"state"  yaml:"state" json:"state"`
	ID             int       `xmlrpc:"id"  yaml:"id" json:"id"`
	EndTime        time.Time `xmlrpc:"end_time"  yaml:"end_time" json:"end_time"`
	SubmitTime     time.Time `xmlrpc:"submit_time"  yaml:"submit_time" json:"submit_time"`
	FailureComment string    `xmlrpc:"failure_comment"  yaml:"failure_comment" json:"failure_comment"`
	Status         int       `xmlrpc:"status"  yaml:"status" json:"status"`
	HealthCheck    bool      `xmlrpc:"health_check"  yaml:"health_check" json:"health_check"`
	Pipeline       bool      `xmlrpc:"pipeline"  yaml:"pipeline" json:"pipeline"`
	Tags           []string  `xmlrpc:"tags"  yaml:"tags" json:"tags"`
	Visibility     string    `xmlrpc:"visibility"  yaml:"visibility" json:"visibility"`
	Submitter      string    `xmlrpc:"submitter"  yaml:"submitter" json:"submitter"`
	StartTime      time.Time `xmlrpc:"start_time"  yaml:"start_time" json:"start_time"`
	Health         string    `xmlrpc:"health"  yaml:"health" json:"health"`
}

func (c Connection) JobsShow(id int) (*JobState, error) {
	var ret JobState

	err := c.con.Call("scheduler.jobs.show", id, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

//TimeoutStruct represents a timeout in a LAVA job definition
type TimeoutStruct struct {
	Seconds int `yaml:"seconds,omitempty"`
	Minutes int `yaml:"minutes,omitempty"`
	Hours   int `yaml:"hours,omitempty"`
}

//TimeoutsStruct contains timeouts defined in a LAVA job definition
type TimeoutsStruct struct {
	Job        TimeoutStruct `yaml:"job,omitempty"`
	Action     TimeoutStruct `yaml:"action,omitempty"`
	Connection TimeoutStruct `yaml:"connection,omitempty"`
}

//ContextStruct contains arch specific in a LAVA job definition
type ContextStruct struct {
	Architecture     string   `yaml:"arch,omitempty"`
	NoKVM            bool     `yaml:"no_kvm"`
	Machine          string   `yaml:"machine,omitempty"`
	CPU              string   `yaml:"cpu,omitempty"`
	GuestFsInterface string   `yaml:"guestfs_interface,omitempty"`
	ExtraOptions     []string `yaml:"extra_options,omitempty"`
}

//ImageStructRamdisk represents a ramdisk(initrd) in a LAVA job definition
type ImageStructRamdisk struct {
	Arguments      string `yaml:"image_arg,omitempty"`
	URL            string `yaml:"url"`
	Compression    string `yaml:"compression,omitempty"`
	Type           string `yaml:"type,omitempty"`
	InstallOverlay bool   `yaml:"install_overlay"`
	InstallModules bool   `yaml:"install_modules"`
}

//ImageStruct represents an arbitrary data image in a LAVA job definition
type ImageStruct struct {
	Arguments   string `yaml:"image_arg,omitempty"`
	URL         string `yaml:"url"`
	Compression string `yaml:"compression,omitempty"`
	Type        string `yaml:"type,omitempty"`
}

//ImagesStruct represents the rootfs and firmware image in a LAVA job definition
type ImagesStruct struct {
	RootFS   ImageStruct `yaml:"rootfs,omitempty"`
	Firmware ImageStruct `yaml:"firmware,omitempty"`
}

//DeployStruct represents the images, OS, timeouts and deploy method in a LAVA job definition
type DeployStruct struct {
	Timeout TimeoutStruct      `yaml:"timeout,omitempty"`
	To      string             `yaml:"to"`
	Images  ImagesStruct       `yaml:"images,omitempty"`
	OS      string             `yaml:"os,omitempty"`
	Kernel  ImageStruct        `yaml:"kernel,omitempty"`
	Ramdisk ImageStructRamdisk `yaml:"ramdisk,omitempty"`
	Dtb     ImageStruct        `yaml:"dtb,omitempty"`
}

//AutoLoginStruct holds the login information in a LAVA job definition
type AutoLoginStruct struct {
	LoginPrompt string `yaml:"login_prompt"`
	UserName    string `yaml:"username"`
}

//TransferOverlayStruct holds the commands to transfer the overlay in a LAVA job definition
type TransferOverlayStruct struct {
	DownloadCommand string `yaml:"download_command"`
	UnpackCommand   string `yaml:"unpack_command"`
}

//BootStruct holds information used to boot a DUT in a LAVA job definition
type BootStruct struct {
	Timeout         TimeoutStruct         `yaml:"timeout,omitempty"`
	Method          string                `yaml:"method"`
	Media           string                `yaml:"media,omitempty"`
	Prompts         []string              `yaml:"prompts,omitempty"`
	FailureMessage  string                `yaml:"failure_message,omitempty"`
	FailureRetry    int                   `yaml:"failure_retry,omitempty"`
	AutoLogin       AutoLoginStruct       `yaml:"auto_login,omitempty"`
	Commands        interface{}           `yaml:"commands,omitempty"`
	TransferOverlay TransferOverlayStruct `yaml:"transfer_overlay,omitempty"`
}

//DefinitionStruct holds the test definition in a LAVA job definition
type DefinitionStruct struct {
	Repository string            `yaml:"repository"`
	From       string            `yaml:"from"`
	Path       string            `yaml:"path"`
	Name       string            `yaml:"name"`
	Parameters map[string]string `yaml:"parameters,omitempty"`
}

//TestStruct holds the test definitions in a LAVA job definition
type TestStruct struct {
	Timeout     TimeoutStruct      `yaml:"timeout,omitempty"`
	Definitions []DefinitionStruct `yaml:"definitions"`
}

//JobStruct represents the job and contains all the structs defined above in a LAVA job definition
type JobStruct struct {
	DeviceType string                   `yaml:"device_type"`
	Context    ContextStruct            `yaml:"context"`
	JobName    string                   `yaml:"job_name"`
	Timeouts   TimeoutsStruct           `yaml:"timeouts"`
	Priority   string                   `yaml:"priority"`
	Visibility string                   `yaml:"visibility"`
	Notify     NotifyStruct             `yaml:"notify,omitempty"`
	Metadata   map[string]string        `yaml:"metadata,omitempty"`
	Actions    []map[string]interface{} `yaml:"actions"`
	Tags       []string                 `yaml:"tags,omitempty"`
}

//CallbackStruct represents the target in a LAVA NotifyStruct
type CallbackStruct struct {
	URL         string `yaml:"url,omitempty"`
	Method      string `yaml:"method,omitempty"`
	Dataset     string `yaml:"dataset,omitempty"`
	ContentType string `yaml:"content-type,omitempty"`
}

//CriteriaStruct represents the status in a LAVA NotifyStruct
type CriteriaStruct struct {
	Status string `yaml:"status,omitempty"`
}

//NotifyStruct represents the notify in a LAVA definition file
type NotifyStruct struct {
	Criteria CriteriaStruct `yaml:"criteria,omitempty"`
	Callback CallbackStruct `yaml:"callback,omitempty"`
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

func (c Connection) JobsSubmitString(def string) ([]int, error) {
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

func (c Connection) JobsSubmit(def *JobStruct) ([]int, error) {

	yaml, err := yaml.Marshal(def)
	if err != nil {
		return []int{-1}, err
	}
	return c.JobsSubmitString(string(yaml))
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
