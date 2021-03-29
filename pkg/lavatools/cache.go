package lavatools

import (
	"github.com/siro20/lavacli/pkg/lava"

	"time"
)

type tagsCache struct {
	Timestamp time.Time
	Tags      []string
}

type deviceCache struct {
	Timestamp time.Time
	Device    lava.Device
}

type deviceListCache struct {
	Timestamp  time.Time
	DeviceList []lava.DeviceList
}

type deviceTypeTemplateCache struct {
	Timestamp time.Time
	Template  string
}

//cache implements a caching layer for lavatools
type cache struct {
	c              *lava.Connection
	pollInterval   time.Duration
	invalidTimeout time.Duration
	retry          *retry
	// Following entries are "cached"
	deviceList         deviceListCache
	devices            map[string]deviceCache
	deviceTags         map[string]tagsCache
	deviceTypeTemplate map[string]deviceTypeTemplateCache
}

func (c *cache) updateDeviceTagsList(name string) (err error) {
	var l []string
	tag, ok := c.deviceTags[name]
	if !ok || time.Now().Sub(tag.Timestamp) > c.pollInterval {
		l, err = c.retry.GetDeviceTagsList(name)
		if err == nil {
			c.deviceTags[name] = tagsCache{Tags: l, Timestamp: time.Now()}
		} else if time.Now().Sub(tag.Timestamp) > c.invalidTimeout {
			return
		}
		err = nil
	}

	return
}

func (c *cache) updateDeviceTypeTemplates(name string) (err error) {
	var s string
	template, ok := c.deviceTypeTemplate[name]
	if !ok || time.Now().Sub(template.Timestamp) > c.pollInterval || len(c.deviceTypeTemplate) == 0 {
		s, err = c.retry.GetDeviceTypeTemplates(name)
		if err == nil {
			c.deviceTypeTemplate[name] = deviceTypeTemplateCache{Template: s, Timestamp: time.Now()}
		} else if time.Now().Sub(template.Timestamp) > c.invalidTimeout {
			return
		}
		err = nil
	}

	return
}

func (c *cache) updateDeviceList() (err error) {
	var l []lava.DeviceList
	if time.Now().Sub(c.deviceList.Timestamp) > c.pollInterval || len(c.deviceList.DeviceList) == 0 {
		l, err = c.retry.GetDeviceList()
		if err == nil {
			c.deviceList.Timestamp = time.Now()
			c.deviceList.DeviceList = l
		} else if time.Now().Sub(c.deviceList.Timestamp) > c.invalidTimeout {
			return
		}
		err = nil
	}

	return
}

func (c *cache) updateDevice(name string) (err error) {
	var d *lava.Device
	dev, ok := c.devices[name]
	if !ok || time.Now().Sub(dev.Timestamp) > c.pollInterval || len(c.devices) == 0 {
		d, err = c.retry.GetDevice(name)
		if err == nil {
			c.devices[name] = deviceCache{Device: *d, Timestamp: time.Now()}
		} else if time.Now().Sub(dev.Timestamp) > c.invalidTimeout {
			return
		}
		err = nil
	}

	return
}

// GetDeviceList returns a cached version of the DeviceList if withing time boundaries
// or returns an error if retrieving new data fails for too long
func (c *cache) GetDeviceList() (devList []lava.DeviceList, err error) {
	err = c.updateDeviceList()
	devList = c.deviceList.DeviceList

	return
}

// GetDevice returns a cached version of the DeviceShow if withing time boundaries
// or returns an error if retrieving new data fails for too long
func (c *cache) GetDevice(name string) (dev lava.Device, err error) {
	err = c.updateDevice(name)
	dev = c.devices[name].Device

	return
}

// GetDeviceTypeTemplates returns a cached version of the DeviceTypeTemplateGet if withing time boundaries
// or returns an error if retrieving new data fails for too long
func (c *cache) GetDeviceTypeTemplates(name string) (template string, err error) {
	err = c.updateDeviceTypeTemplates(name)
	template = c.deviceTypeTemplate[name].Template

	return
}

// GetDeviceTagsList returns a cached version of the DeviceTagList if withing time boundaries
// or returns an error if retrieving new data fails for too long
func (c *cache) GetDeviceTagsList(name string) (tags []string, err error) {
	err = c.updateDeviceTagsList(name)
	tags = c.deviceTags[name].Tags

	return
}

// updatePeriodic updates the cache on a periodic interval
func (c *cache) updatePeriodic(timeout time.Duration) {
	time.Sleep(timeout)
	time.Sleep(1)
	c.updateDeviceList()
	for i := range c.deviceList.DeviceList {
		c.updateDeviceTagsList(c.deviceList.DeviceList[i].Hostname)
		c.updateDevice(c.deviceList.DeviceList[i].Hostname)
	}
	go c.updatePeriodic(timeout)
}

//newLavaToolsCache returns a cache object
func newLavaToolsCache(c *lava.Connection, retry *retry, opt Options) (obj *cache, err error) {
	obj = &cache{c: c,
		deviceList:         deviceListCache{DeviceList: []lava.DeviceList{}},
		devices:            map[string]deviceCache{},
		deviceTags:         map[string]tagsCache{},
		deviceTypeTemplate: map[string]deviceTypeTemplateCache{},
		pollInterval:       opt.PollInterval,
		invalidTimeout:     opt.InvalidTimeout,
		retry:              retry,
	}

	if opt.BackgroundPrefetching {
		go obj.updatePeriodic(opt.BackgroundInterval)
	}

	return
}
