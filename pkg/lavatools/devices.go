package lavatools

import (
	"fmt"
	"strings"
	"time"

	"github.com/siro20/lavacli/pkg/lava"
)

// DeviceListWithRetry returns the device list
// Retries to get the list in case of error
func (con lt) DeviceListWithRetry() (devList []lava.DeviceList, err error) {
	for i := 0; i < 5; i++ {
		devList, err = con.c.DevicesList()
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}

	return
}

//DevicesTagsListWithRetry returns the tags of the device specified by name
// Retries to get the list in case of error
func (con lt) DevicesTagsListWithRetry(name string) (ret []string, err error) {
	for i := 0; i < 5; i++ {
		ret, err = con.c.DevicesTagsList(name)
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		break
	}

	return
}

func (con lt) updateDevice(name string) (err error) {
	var d *lava.Device
	dev, ok := con.devices[name]
	if !ok || time.Now().Sub(dev.Timestamp) > con.pollInterval || len(con.devices) == 0 {
		d, err = con.c.DevicesShow(name)
		if err == nil {
			con.devices[name] = deviceCache{Device: *d, Timestamp: time.Now()}
		} else if time.Now().Sub(dev.Timestamp) > con.invalidTimeout {
			return
		}
		err = nil
	}

	return
}

//DevicesShowCached caches the device to show as every API call takes a while
func (con lt) DevicesShowCached(name string) (dev lava.Device, err error) {
	var ok bool
	var d deviceCache

	err = con.updateDevice(name)
	d, ok = con.devices[name]
	if !ok {
		err = fmt.Errorf("%s not found in cache", name)
	}
	dev = d.Device

	return
}

func (con *lt) updateDeviceList() (err error) {
	var l []lava.DeviceList
	if time.Now().Sub(con.deviceList.Timestamp) > con.pollInterval || len(con.deviceList.DeviceList) == 0 {
		l, err = con.DeviceListWithRetry()
		if err == nil {
			con.deviceList.Timestamp = time.Now()
			con.deviceList.DeviceList = l
		} else if time.Now().Sub(con.deviceList.Timestamp) > con.invalidTimeout {
			return
		}
		err = nil
	}

	return
}

//DeviceListCached caches the device list as every API call takes a while
func (con lt) DeviceListCached() (devList []lava.DeviceList, err error) {

	err = con.updateDeviceList()
	devList = con.deviceList.DeviceList
	return
}

//DeviceListHealthyCached caches the device list of healthy devices as every API call takes a while
func (con lt) DeviceListHealthyCached() (devList []lava.DeviceList, err error) {
	devList = []lava.DeviceList{}

	err = con.updateDeviceList()
	for _, d := range con.deviceList.DeviceList {
		if strings.ToLower(d.Health) == "good" {
			devList = append(devList, d)
		}
	}

	return
}

func (con *lt) updateDeviceTagsList(name string) (err error) {
	var l []string
	tag, ok := con.deviceTags[name]
	if !ok || time.Now().Sub(tag.Timestamp) > con.pollInterval {
		l, err = con.DevicesTagsListWithRetry(name)
		if err == nil {
			con.deviceTags[name] = tagsCache{Tags: l, Timestamp: time.Now()}
		} else if time.Now().Sub(tag.Timestamp) > con.invalidTimeout {
			return
		}
		err = nil
	}

	return
}

//DevicesTagsListCached caches the device tags as every API call takes a while
func (con lt) DevicesTagsListCached(name string) (tags []string, err error) {
	var ok bool
	var t tagsCache
	err = con.updateDeviceTagsList(name)
	t, ok = con.deviceTags[name]
	if !ok {
		err = fmt.Errorf("%s not found in cache", name)
	}
	tags = t.Tags

	return
}

//DeviceHasTagCached uses the cached device tag list to determine if a device has
//the specified tag, with optional case insensitive comparison
func (con lt) DeviceHasTagCached(name string, tag string, ignoreCase bool) (has bool, err error) {
	var deviceTags []string
	deviceTags, err = con.DevicesTagsListCached(name)
	if err != nil {
		return
	}
	for _, t := range deviceTags {
		if ignoreCase {
			if strings.ToLower(tag) == strings.ToLower(t) {
				has = true
				return
			}
		} else {
			if tag == t {
				has = true
				return
			}
		}
	}

	return
}

//DeviceHasTagsCached uses the cached device tag list to determine if a device has
//the specified tags, with optional case insensitive comparison
func (con lt) DeviceHasTagsCached(name string, tags []string, ignoreCase bool) (has bool, err error) {
	for _, t := range tags {
		has, err = con.DeviceHasTagCached(name, t, ignoreCase)
		if !has || err != nil {
			return
		}
	}

	return
}

//DeviceIsAliveCached uses the cached device list to determine if a device is online
//and available for testing
func (con lt) DeviceIsAliveCached(name string) (alive bool, err error) {
	var list []lava.DeviceList
	list, err = con.DeviceListCached()
	if err != nil {
		return
	}
	for _, d := range list {
		if d.Hostname != name {
			continue
		}
		if strings.ToLower(d.Health) == "good" {
			alive = true
			return
		}
	}

	return
}

//DeviceOfTypeIsAliveAndHasTagCached uses the cached device list to find a device that
//is of good health, of the specified type and has all tags set
func (con lt) DeviceOfTypeIsAliveAndHasTagCached(deviceType string, tagsToMatch []string) (alive bool, err error) {
	var list []lava.DeviceList
	list, err = con.DeviceListCached()
	if err != nil {
		return
	}
	for _, d := range list {
		var ok bool
		if strings.ToLower(d.Health) != "good" {
			continue
		}
		if strings.ToLower(d.Type) != strings.ToLower(deviceType) {
			continue
		}
		ok, err = con.DeviceHasTagsCached(d.Hostname, tagsToMatch, false)
		if err != nil {
			break
		}
		if !ok {
			continue
		}

		alive = true
		return
	}

	return
}
