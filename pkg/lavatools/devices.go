package lavatools

import (
	"strings"

	"github.com/siro20/lavacli/pkg/lava"
)

// DeviceListWithRetry returns the device list
// Retries to get the list in case of error
func (con lt) DeviceListWithRetry() (devList []lava.DeviceList, err error) {

	devList, err = con.retry.GetDeviceList()

	return
}

//DevicesTagsListWithRetry returns the tags of the device specified by name
// Retries to get the list in case of error
func (con lt) DevicesTagsListWithRetry(name string) (ret []string, err error) {

	ret, err = con.retry.GetDeviceTagsList(name)

	return
}

//DevicesShowCached caches the device to show as every API call takes a while
func (con lt) DevicesShowCached(name string) (dev lava.Device, err error) {

	dev, err = con.cache.GetDevice(name)

	return
}

//DeviceListCached caches the device list as every API call takes a while
func (con lt) DeviceListCached() (devList []lava.DeviceList, err error) {

	devList, err = con.cache.GetDeviceList()

	return
}

//DeviceListHealthyCached caches the device list of healthy devices as every API call takes a while
func (con lt) DeviceListHealthyCached() (devListHealthy []lava.DeviceList, err error) {
	var devList []lava.DeviceList
	devListHealthy = []lava.DeviceList{}

	devList, err = con.cache.GetDeviceList()
	if err != nil {
		return
	}
	for _, d := range devList {
		if strings.ToLower(d.Health) == "good" {
			devListHealthy = append(devListHealthy, d)
		}
	}

	return
}

//DevicesTagsListCached caches the device tags as every API call takes a while
func (con lt) DevicesTagsListCached(name string) (tags []string, err error) {

	tags, err = con.cache.GetDeviceTagsList(name)

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
