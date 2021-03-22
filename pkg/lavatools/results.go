package lavatools

import (

	"github.com/siro20/lavacli/pkg/lava"
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"
)

//GetJobTestResultsWithRetry returns the results for a given job ID
// Retries to get the result in case of error
func (con lt) GetJobTestResultsWithRetry(id int) (ret lava.Result, err error) {
	var r lava.Result
	for i := 0; i < 5; i++ {
		r, err = con.c.Results(id)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}
	if err != nil {
		return
	}

	// HACK: convert to yaml and back to convert from lava.LavaResult to api.LavaResult
	var d []byte
	d, err = yaml.Marshal(&r)
	if err != nil {
		err = fmt.Errorf("Failed to marshal data: %v", err)
		return
	}

	err = yaml.Unmarshal(d, &ret)
	if err != nil {
		err = fmt.Errorf("Failed to unmarshal data: %v", err)
		return
	}

	return
}
