package lavatools

import (
	"github.com/siro20/lavacli/pkg/lava"

	"time"
)

//rety implements a retry layer for lavatools
type retry struct {
	c          *lava.Connection
	retryCount int
}

func (r *retry) GetDeviceList() (devList []lava.DeviceList, err error) {
	for i := 0; i <= r.retryCount; i++ {
		devList, err = r.c.DevicesList()
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}
	return
}

func (r *retry) GetDevice(name string) (dev *lava.Device, err error) {
	for i := 0; i <= r.retryCount; i++ {
		dev, err = r.c.DevicesShow(name)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}

	return
}

func (r *retry) GetDeviceTypeTemplates(name string) (template string, err error) {
	for i := 0; i <= r.retryCount; i++ {
		template, err = r.c.DevicesTypesTemplateGet(name)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}

		break
	}

	return
}

func (r *retry) GetDeviceTagsList(name string) (tags []string, err error) {
	for i := 0; i <= r.retryCount; i++ {
		tags, err = r.c.DevicesTagsList(name)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}

	return
}

//newLavaToolsRetry returns a retry object
func newLavaToolsRetry(c *lava.Connection, opt Options) (obj *retry, err error) {
	obj = &retry{c: c,
		retryCount: opt.RetryCount,
	}

	return
}
