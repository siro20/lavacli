// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kolo/xmlrpc"
)

// ConnectionOptions allows to pass additional parameters
type ConnectionOptions struct {
	Transport *http.Transport
}

// DefaultOptions must be passed as argument to the Connect.. methods if no overwrites are made
var DefaultOptions = ConnectionOptions{
	Transport: &http.Transport{},
}

// Connection holds metadata used to communicate with the LAVA XMLRPC server
type Connection struct {
	con   *xmlrpc.Client
	proxy string
	uri   string
	opt   ConnectionOptions
}

// ConnectByURI connects to an LAVA XMLRPC server using the provided URI, proxy and transport
func ConnectByURI(uri string, proxy string, opt ConnectionOptions) (*Connection, error) {
	var ret Connection
	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		opt.Transport.Proxy = http.ProxyURL(u)
	}

	client, err := xmlrpc.NewClient(uri, opt.Transport)
	if err != nil {
		return nil, err
	}
	ret.con = client
	ret.proxy = proxy
	ret.uri = uri
	ret.opt = opt

	return &ret, nil
}

// ConnectByConfigID connects to an LAVA XMLRPC server using the provided identity and lavacli.yaml
func ConnectByConfigID(identityName string, opt ConnectionOptions) (*Connection, error) {
	var u string
	var c ConfigIndentity

	configs, err := GetConf()
	if err != nil {
		return nil, err
	}
	found := false
	for k, v := range configs {
		if k == identityName {
			c = v
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("Identity %s not found in lavacli.yaml", identityName)
	}

	if c.URI == "" {
		return nil, fmt.Errorf("No URI found in config")
	}

	if c.Username != "" && c.Token != "" {
		url, err := url.Parse(c.URI)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse URI: %v", err)
		}
		u = fmt.Sprintf("%s://%s:%s@%s%s", url.Scheme, c.Username, c.Token, url.Host, url.Path)
	} else if c.URI != "" {
		u = c.URI
	}

	return ConnectByURI(u, c.Proxy, opt)
}

// ConnectByCredentials connects to an LAVA XMLRPC server using the provided credentials
func ConnectByCredentials(uri string, username string, token string, proxy string, opt ConnectionOptions) (*Connection, error) {
	var u string

	if uri == "" {
		return nil, fmt.Errorf("No URI specified")
	}

	if username != "" && token != "" {
		url, err := url.Parse(uri)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse URI: %v", err)
		}
		u = fmt.Sprintf("%s://%s:%s@%s%s", url.Scheme, username, token, url.Host, url.Path)
	} else if uri != "" {
		u = uri
	}

	return ConnectByURI(u, proxy, opt)
}
