package lavatools

import (
	"github.com/siro20/lavacli/pkg/lava"

	"time"
)

const (
	//JobDefaultTimeout is the maximum duration for a test job
	JobDefaultTimeout = 40 // 20 Minutes for a complete harddisk rewrite on HP8200
	//DefaultTimeout is the maximum duration for a test step
	DefaultTimeout     = 20
	QEMUDefaultTimeout = 30
	defaultPriority    = "medium"
	defaultVisibility  = "public"
	// DefaultKexecDeployTimeout is the time in minutes to wait for kexec a kernel
	DefaultKexecDeployTimeout = 10
	// DefaultHarddiskDeployTimeout is the time in minutes to wait for a write to harddisk to finish
	DefaultHarddiskDeployTimeout = 30
	// DefaultFlashDeployTimeout is the time in minutes to wait for a write to flash to finish
	DefaultFlashDeployTimeout = 10
)

// Timeouts are only used within jobs steps if the timeout field is missing
type Timeouts struct {
	// Job is the maximum duration for a test job
	Job int
	// Step is the maximum duration for a test step
	Step int
	// QEMU is the maximum duration for a test step in qemu
	QEMU int
	// KexecDeploy is the time in minutes to wait for kexec a kernel
	KexecDeploy int
	// HarddiskDeploy is the time in minutes to wait for a write to harddisk to finish
	HarddiskDeploy int
	// FlashDeploy is the time in minutes to wait for a write to flash to finish
	FlashDeploy int
}

// JobOptions contain additional settings for the job submission
type JobOptions struct {
	Timeout    Timeouts
	Priority   string
	Visibility string
}

// DefaultJobOptions are to be used as default
var DefaultJobOptions = JobOptions{
	Timeout: Timeouts{
		Job:            JobDefaultTimeout,
		Step:           DefaultTimeout,
		QEMU:           QEMUDefaultTimeout,
		KexecDeploy:    DefaultKexecDeployTimeout,
		HarddiskDeploy: DefaultHarddiskDeployTimeout,
		FlashDeploy:    DefaultFlashDeployTimeout,
	},
	Priority:   defaultPriority,
	Visibility: defaultVisibility,
}

type Options struct {
	//RetryCount is the number of retries done in the *Retry methods
	//before the error is returned to the caller
	RetryCount int
	//PollInterval sets the interval between two data fetches
	//Errors fetching new data are not returned as long as
	//InvalidTimeout isn't reached.
	PollInterval time.Duration
	//InvalidTimeout sets the time after cached data is discarded and
	//errors fetching new data are returned again
	InvalidTimeout time.Duration
	//BackgroundPrefetching enables data fetching in a go routine for all
	//*Cached methods. This cannot be canceled until now
	BackgroundPrefetching bool
	//BackgroundInterval sets the interval between fetching data
	//in a go routing
	BackgroundInterval time.Duration
}

var DefaultOptions = Options{
	RetryCount:            5,
	PollInterval:          time.Minute * 5,
	InvalidTimeout:        time.Minute * 10,
	BackgroundPrefetching: true,
	BackgroundInterval:    time.Minute * 5,
}

//lt implements the Lavatools interface
type lt struct {
	c              *lava.Connection
	pollInterval   time.Duration
	invalidTimeout time.Duration
	retryCount     int
	cache          *cache
	retry          *retry
}

type Lavatools interface {
	// jobs
	LaunchJob(job lava.JobStruct, opt JobOptions) (int, error)
	JobsValidate(jobYaml string) (msg string, err error)
	JobsSubmitWithRetry(job *lava.JobStruct) (id int, err error)
	JobsSubmitStringWithRetry(jobYaml string) (id int, err error)
	JobsShowWithRetry(id int) (state *lava.JobState, err error)
	JobsDefinitionWithRetry(id int) (job *lava.JobStruct, err error)
	QueryJobListWithRetry(state string, health string, start int, limit int) (list []lava.JobsListing, err error)
	CancelJobWithRetry(id int) (err error)
	// results
	GetJobTestResultsWithRetry(id int) (ret lava.Result, err error)
	// device
	DeviceListWithRetry() (defList []lava.DeviceList, err error)
	// device tags
	DevicesTagsListWithRetry(name string) (ret []string, err error)
	// device-types template
	DevicesTypesTemplateGetWithRetry(name string) (ret string, err error)
	// following functions use cached results
	DeviceListCached() (devList []lava.DeviceList, err error)
	DeviceListHealthyCached() (devList []lava.DeviceList, err error)
	DevicesTypesTemplateGetCached(name string) (ret string, err error)
	DeviceHasTagCached(name string, tag string, ignoreCase bool) (has bool, err error)
	DeviceHasTagsCached(name string, tags []string, ignoreCase bool) (has bool, err error)
	DeviceOfTypeIsAliveAndHasTagCached(deviceType string, tagsToMatch []string) (alive bool, err error)
}

//NewLavaTools returns an interface to Lavatools
func NewLavaTools(c *lava.Connection, opt Options) (con Lavatools, err error) {
	retry, err := newLavaToolsRetry(c, opt)
	if err != nil {
		return
	}

	cache, err := newLavaToolsCache(c, retry, opt)
	if err != nil {
		return
	}

	obj := lt{c: c,
		pollInterval:   opt.PollInterval,
		invalidTimeout: opt.InvalidTimeout,
		retryCount:     opt.RetryCount,
		cache:          cache,
		retry:          retry,
	}
	con = obj
	return
}
