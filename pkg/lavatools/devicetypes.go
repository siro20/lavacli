package lavatools

import (


	"time"
	"fmt"
)

// DevicesTypesTemplateGetWithRetry returns the LAVA type specific 'template'
// Retries to get the template in case of error
func (con lt) DevicesTypesTemplateGetWithRetry(name string) (ret string, err error) {
	for i := 0; i < 5; i++ {
		ret, err = con.c.DevicesTypesTemplateGet(name)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}

	return
}

// DevicesTypesTemplateGetCached returns the cached device type template 
func (con lt) DevicesTypesTemplateGetCached(name string) (ret string, err error) {

	err = con.updateDeviceTypeTemplates(name)
	if err != nil {
		return
	}
	t, ok := con.deviceTypeTemplate[name]
	if !ok {
		err = fmt.Errorf("%s not found in cache", name)
	}
	ret = t.Template
	return
}

func (con *lt) updateDeviceTypeTemplates(name string) (err error) {
	var s string
	template, ok := con.deviceTypeTemplate[name]
	if !ok || time.Now().Sub(template.Timestamp) > con.pollInterval || len(con.deviceTypeTemplate) == 0 {
			s, err = con.DevicesTypesTemplateGetWithRetry(name)
			if err == nil {
				con.deviceTypeTemplate[name] = deviceTypeTemplateCache{Template:s, Timestamp: time.Now()}
			} else if time.Now().Sub(template.Timestamp) > con.invalidTimeout {
				return
			}
			err = nil
	}

	return
}