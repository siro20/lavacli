package lavatools

// DevicesTypesTemplateGetWithRetry returns the LAVA type specific 'template'
// Retries to get the template in case of error
func (con lt) DevicesTypesTemplateGetWithRetry(name string) (ret string, err error) {

	ret, err = con.retry.GetDeviceTypeTemplates(name)

	return
}

// DevicesTypesTemplateGetCached returns the cached device type template
func (con lt) DevicesTypesTemplateGetCached(name string) (ret string, err error) {

	ret, err = con.cache.GetDeviceTypeTemplates(name)
	return
}
