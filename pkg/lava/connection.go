// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kolo/xmlrpc"
)

type LavaConnection struct {
	Con   *xmlrpc.Client
	proxy string
	uri   string
}

func LavaConnectByUri(uri string, proxy string) (*LavaConnection, error) {
	var ret LavaConnection
	tr := http.Transport{}

	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		tr.Proxy = http.ProxyURL(u)
	}

	client, err := xmlrpc.NewClient(uri, &tr)
	if err != nil {
		return nil, err
	}
	ret.Con = client
	ret.proxy = proxy
	ret.uri = uri

	return &ret, nil
}

func LavaConnectByConfigID(identityName string) (*LavaConnection, error) {
	var u string
	var c LavaConfigIndentity

	configs, err := LavaGetConf()
	if err != nil {
		return nil, err
	}
	found := false
	for k, v := range configs {
		if k == identityName {
			c = v
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Identity %s not found in lavacli.yaml\n", identityName)
	}

	if c.Uri == "" {
		return nil, fmt.Errorf("No URI found in config\n")
	}

	if c.Username != "" && c.Token != "" {
		url, err := url.Parse(c.Uri)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse URI: %v\n", err)
		}
		u = fmt.Sprintf("%s://%s:%s@%s%s", url.Scheme, c.Username, c.Token, url.Host, url.Path)
	} else if c.Uri != "" {
		u = c.Uri
	}

	return LavaConnectByUri(u, c.Proxy)
}

func LavaConnectByCredentials(uri string, username string, token string, proxy string) (*LavaConnection, error) {
	var u string

	if uri == "" {
		return nil, fmt.Errorf("No URI specified\n")
	}

	if username != "" && token != "" {
		url, err := url.Parse(uri)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse URI: %v\n", err)
		}
		u = fmt.Sprintf("%s://%s:%s@%s%s", url.Scheme, username, token, url.Host, url.Path)
	} else if uri != "" {
		u = uri
	}

	return LavaConnectByUri(u, proxy)
}
