package lavatools

import (
	"testing"
	"time"

	"github.com/siro20/lavacli/pkg/lava"
)

func Test_lt_DeviceListHealthyCached(t *testing.T) {

	lavacon, err := lava.ConnectByConfigID("", lava.DefaultOptions)
	if err != nil {
		t.Skip()
	}
	con, err := NewLavaTools(lavacon, DefaultOptions)
	if err != nil {
		t.Skip()
	}

	startUncached := time.Now()
	_, err = lavacon.DevicesList()
	if err != nil {
		t.Errorf("DeviceListHealthyCached() got unexpected error = %v", err)
		return
	}

	elapsedUncached := time.Since(startUncached)

	_, err = con.DeviceListHealthyCached()
	if err != nil {
		t.Errorf("DeviceListHealthyCached() got unexpected error = %v", err)
		return
	}

	startCached := time.Now()
	_, err = con.DeviceListHealthyCached()
	if err != nil {
		t.Errorf("DeviceListHealthyCached() got unexpected error = %v", err)
		return
	}
	elapsedCached := time.Since(startCached)
	if elapsedCached*100 > elapsedUncached {
		t.Errorf("DeviceListHealthyCached took very long to return: cached = %v, uncached = %v", elapsedCached, elapsedUncached)
		return
	}
}

func Test_lt_DeviceListCached(t *testing.T) {

	lavacon, err := lava.ConnectByConfigID("", lava.DefaultOptions)
	if err != nil {
		t.Skip()
	}
	con, err := NewLavaTools(lavacon, DefaultOptions)
	if err != nil {
		t.Skip()
	}

	startUncached := time.Now()
	_, err = lavacon.DevicesList()
	if err != nil {
		t.Errorf("DeviceListCached() got unexpected error = %v", err)
		return
	}

	elapsedUncached := time.Since(startUncached)

	_, err = con.DeviceListCached()
	if err != nil {
		t.Errorf("DeviceListCached() got unexpected error = %v", err)
		return
	}

	startCached := time.Now()
	_, err = con.DeviceListCached()
	if err != nil {
		t.Errorf("DeviceListCached() got unexpected error = %v", err)
		return
	}
	elapsedCached := time.Since(startCached)
	if elapsedCached*100 > elapsedUncached {
		t.Errorf("DeviceListCached took very long to return: cached = %v, uncached = %v", elapsedCached, elapsedUncached)
		return
	}
}

func Test_lt_DevicesTypesTemplateGetCached(t *testing.T) {

	lavacon, err := lava.ConnectByConfigID("", lava.DefaultOptions)
	if err != nil {
		t.Skip()
	}
	con, err := NewLavaTools(lavacon, DefaultOptions)
	if err != nil {
		t.Skip()
	}

	startUncached := time.Now()
	devTypes, err := lavacon.DevicesTypesList(false)
	if err != nil {
		t.Errorf("DeviceListCached() got unexpected error = %v", err)
		return
	}

	elapsedUncached := time.Since(startUncached)

	for _, dtype := range devTypes {
		_, err = con.DevicesTypesTemplateGetCached(dtype.Name)
		if err != nil {
			t.Errorf("DevicesTypesTemplateGetCached(%s) got unexpected error = %v", dtype.Name, err)
			return
		}

		startCached := time.Now()
		_, err = con.DevicesTypesTemplateGetCached(dtype.Name)
		if err != nil {
			t.Errorf("DevicesTypesTemplateGetCached(%s) got unexpected error = %v", dtype.Name, err)
			return
		}
		elapsedCached := time.Since(startCached)
		if elapsedCached*100 > elapsedUncached {
			t.Errorf("DevicesTypesTemplateGetCached(%s) took very long to return: cached = %v, uncached = %v", dtype.Name, elapsedCached, elapsedUncached)
			return
		}
	}
}
