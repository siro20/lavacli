package lavatools

import (
	"fmt"
	"log"
	"time"

	"github.com/siro20/lavacli/pkg/lava"

	yaml "gopkg.in/yaml.v2"
)

//updateTimeouts sets the default timeouts and accumulates timeouts for the global job timeout
func updateTimeouts(job *lava.JobStruct, opt JobOptions) error {
	// Fill default timeouts

	for i := range job.Actions {
		if deploy, ok := job.Actions[i]["deploy"]; ok {
			deploy, ok := deploy.(lava.DeployStruct)
			if !ok {
				continue
			}
			if deploy.Timeout.Seconds == 0 && deploy.Timeout.Minutes == 0 && deploy.Timeout.Hours == 0 {
				if deploy.To == "flasher" {
					deploy.Timeout.Minutes = opt.Timeout.FlashDeploy
				} else if deploy.To == "tftp" {
					if deploy.Ramdisk.URL != "" {
						deploy.Timeout.Minutes = opt.Timeout.KexecDeploy
					} else {
						deploy.Timeout.Minutes = opt.Timeout.HarddiskDeploy
					}
				}
				job.Actions[i]["deploy"] = deploy
			}
		}
		if boot, ok := job.Actions[i]["boot"]; ok {
			boot, ok := boot.(lava.BootStruct)
			if !ok {
				continue
			}
			if boot.Timeout.Seconds == 0 && boot.Timeout.Minutes == 0 && boot.Timeout.Hours == 0 {
				if boot.Method == "u-boot" {
					boot.Timeout.Minutes = opt.Timeout.HarddiskDeploy
				} else if boot.Method == "qemu" {
					boot.Timeout.Minutes = opt.Timeout.QEMU
				} else if boot.AutoLogin.LoginPrompt != "" {
					boot.Timeout.Minutes = opt.Timeout.Step
				}

				job.Actions[i]["boot"] = boot
			}
		}
		if test, ok := job.Actions[i]["test"]; ok {
			test, ok := test.(lava.TestStruct)

			if !ok {
				continue
			}

			if test.Timeout.Seconds == 0 && test.Timeout.Minutes == 0 && test.Timeout.Hours == 0 {
				test.Timeout.Minutes = opt.Timeout.Step
				job.Actions[i]["test"] = test
			}
		}
	}

	actionDeployTimeout := 0
	actionBootTimeout := 0
	actionTestTimeout := 0
	actiontimeout := 0

	for _, action := range job.Actions {
		if deploy, ok := action["deploy"]; ok {
			deploy, ok := deploy.(lava.DeployStruct)
			if !ok {
				continue
			}
			if deploy.Timeout.Seconds > 0 {
				actionDeployTimeout += deploy.Timeout.Seconds
			}
			if deploy.Timeout.Minutes*60 > 0 {
				actionDeployTimeout += deploy.Timeout.Minutes * 60
			}
			if deploy.Timeout.Hours*3600 > 0 {
				actionDeployTimeout += deploy.Timeout.Hours * 3600
			}
			if deploy.Timeout.Seconds > actiontimeout {
				actiontimeout = deploy.Timeout.Seconds
			}
			if deploy.Timeout.Minutes*60 > actiontimeout {
				actiontimeout = deploy.Timeout.Minutes * 60
			}
			if deploy.Timeout.Hours*3600 > actiontimeout {
				actiontimeout = deploy.Timeout.Hours * 3600
			}
		}
		if boot, ok := action["boot"]; ok {
			boot, ok := boot.(lava.BootStruct)
			if !ok {
				continue
			}
			if boot.Timeout.Seconds > 0 {
				actionBootTimeout += boot.Timeout.Seconds
			}
			if boot.Timeout.Minutes*60 > 0 {
				actionBootTimeout += boot.Timeout.Minutes * 60
			}
			if boot.Timeout.Hours*3600 > 0 {
				actionBootTimeout += boot.Timeout.Hours * 3600
			}
			if boot.Timeout.Seconds > actiontimeout {
				actiontimeout = boot.Timeout.Seconds
			}
			if boot.Timeout.Minutes*60 > actiontimeout {
				actiontimeout = boot.Timeout.Minutes * 60
			}
			if boot.Timeout.Hours*3600 > actiontimeout {
				actiontimeout = boot.Timeout.Hours * 3600
			}
		}
		if test, ok := action["test"]; ok {
			test, ok := test.(lava.TestStruct)

			if !ok {
				continue
			}
			if test.Timeout.Seconds > 0 {
				actionTestTimeout += test.Timeout.Seconds
			}
			if test.Timeout.Minutes*60 > 0 {
				actionTestTimeout += test.Timeout.Minutes * 60
			}
			if test.Timeout.Hours*3600 > 0 {
				actionTestTimeout += test.Timeout.Hours * 3600
			}
			if test.Timeout.Seconds > actiontimeout {
				actiontimeout = test.Timeout.Seconds
			}
			if test.Timeout.Minutes*60 > actiontimeout {
				actiontimeout = test.Timeout.Minutes * 60
			}
			if test.Timeout.Hours*3600 > actiontimeout {
				actiontimeout = test.Timeout.Hours * 3600
			}
		}
	}

	timeout := JobDefaultTimeout * 60

	if actionDeployTimeout+actionBootTimeout+actionTestTimeout > timeout {
		timeout = actionDeployTimeout + actionBootTimeout + actionTestTimeout
	}

	job.Timeouts = lava.TimeoutsStruct{
		Job:        lava.TimeoutStruct{Seconds: timeout},
		Action:     lava.TimeoutStruct{Seconds: actiontimeout},
		Connection: lava.TimeoutStruct{Seconds: timeout},
	}

	return nil
}

//LaunchJob starts a new job on the LAVA master server
func (con lt) LaunchJob(job lava.JobStruct, opt JobOptions) (int, error) {
	// Fill data

	err := updateTimeouts(&job, opt)
	if err != nil {
		return -1, err
	}
	job.Priority = opt.Priority
	job.Visibility = opt.Visibility

	yaml, err := yaml.Marshal(&job)
	if err != nil {
		return -1, err
	}

	// Some basic checks before submission
	output, err := con.JobsValidate(string(yaml))
	if err != nil {
		log.Printf("Invalid YAML job definition: %s, %s\n", output, err.Error())
		return -1, err
	}

	jobID, err := con.JobsSubmitStringWithRetry(string(yaml))
	if err != nil {
		return -1, err
	}

	return jobID, nil
}

//JobsValidate validates the job definition
func (con lt) JobsValidate(jobYaml string) (msg string, err error) {

	jobErrors, err := con.c.JobsValidate(jobYaml, false)
	if err != nil {
		return
	}

	for k, v := range jobErrors {
		msg += fmt.Sprintf("Error: %s %v\n", k, v)
	}

	return
}

// JobsSubmitWithRetry starts a new job with the definition specified as string and returns the job ID
// Retries to get the definition in case of error
func (con lt) JobsSubmitWithRetry(job *lava.JobStruct) (id int, err error) {
	var ids []int
	id = -1
	ids, err = con.c.JobsSubmit(job)
	if err != nil {
		err = fmt.Errorf("JobsSubmit returned error: %v\n", err)
		return
	}
	if len(ids) == 0 {
		err = fmt.Errorf("Got Invalid jobIDs")
		return
	}
	id = ids[0]
	return
}

// JobsSubmitStringWithRetry starts a new job with the definition specified as string and returns the job ID
// Retries to get the definition in case of error
func (con lt) JobsSubmitStringWithRetry(jobYaml string) (id int, err error) {
	var ids []int
	id = -1
	ids, err = con.c.JobsSubmitString(jobYaml)
	if err != nil {
		err = fmt.Errorf("JobsSubmit returned error: %v\n", err)
		return
	}
	if len(ids) == 0 {
		err = fmt.Errorf("Got Invalid jobIDs")
		return
	}
	id = ids[0]
	return
}

//JobsShowWithRetry returns the job state for a given job ID
func (con lt) JobsShowWithRetry(id int) (state *lava.JobState, err error) {

	for i := 0; i < 5; i++ {
		state, err = con.c.JobsShow(id)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}

	return
}

//JobsDefinitionWithRetry returns the job definition for a given job ID
// Retries to get the definition in case of error
func (con lt) JobsDefinitionWithRetry(id int) (job *lava.JobStruct, err error) {
	var def lava.JobDefintion
	for i := 0; i < 5; i++ {
		def, err = con.c.JobsDefinition(id)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}
	if err != nil {
		return
	}

	err = yaml.Unmarshal([]byte(def), &job)
	if err != nil {
		err = fmt.Errorf("Failed to unmarshal job definition: %s. The input was '%s'", err.Error(), string(def))
		return
	}

	return
}

//CancelJobWithRetry cancels the job of a given job ID
func (con lt) CancelJobWithRetry(id int) (err error) {
	for i := 0; i < 5; i++ {
		err = con.c.JobsCancel(id)

		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}
	return
}

//QueryJobListWithRetry returns a list of jobs matching the specified filters
// Retries to get the list in case of error
func (con lt) QueryJobListWithRetry(state string, health string, start int, limit int) (list []lava.JobsListing, err error) {
	for i := 0; i < 5; i++ {
		list, err = con.c.JobsList(state, health, start, limit)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}

	return
}
