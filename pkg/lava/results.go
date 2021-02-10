// SPDX-License-Identifier: BSD-3-Clause

package lava

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type LavaResult []struct {
	Name         string `yaml:"name,omitempty" json:"name,omitempty"`
	Result       string `yaml:"result,omitempty" json:"result,omitempty"`
	Id           string `yaml:"id,omitempty" json:"id,omitempty"`
	Job          string `yaml:"job,omitempty" json:"job,omitempty"`
	Level        string `yaml:"level,omitempty" json:"level,omitempty"`
	LogLineEnd   string `yaml:"log_end_line,omitempty" json:"log_end_line,omitempty"`
	LogLineStart string `yaml:"log_start_line,omitempty" json:"log_start_line,omitempty"`
	Suite        string `yaml:"suite,omitempty" json:"suite,omitempty"`
	URL          string `yaml:"url,omitempty" json:"url,omitempty"`
	Measurement  string `yaml:"measurement,omitempty" json:"measurement,omitempty"`
	Logged       string `yaml:"logged,omitempty" json:"logged,omitempty"`
	Metadata     struct {
		Case       string `yaml:"case,omitempty" json:"case,omitempty"`
		Definition string `yaml:"definition,omitempty" json:"definition,omitempty"`
		Duration   string `yaml:"duration,omitempty" json:"duration,omitempty"`
		// This is broken as it should contain a slice of maps with a single key,
		// but it can be anything...
		//Extra      []map[string]interface{} `yaml:"extra,omitempty"`
		Level      string `yaml:"level,omitempty" json:"level,omitempty"`
		Namespace  string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
		Result     string `yaml:"result,omitempty" json:"result,omitempty"`
		UUID       string `yaml:"uuid,omitempty" json:"uuid,omitempty"`
		Revision   string `yaml:"revision,omitempty" json:"revision,omitempty"`
		Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
		CommitID   string `yaml:"commit_id,omitempty" json:"commit_id,omitempty"`
		Path       string `yaml:"path,omitempty" json:"path,omitempty"`
		ErrorMsg   string `yaml:"error_msg,omitempty" json:"error_msg,omitempty"`
		ErrorType  string `yaml:"error_type,omitempty" json:"error_type,omitempty"`
	} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

func (c LavaConnection) LavaResultsAsYAML(id int) (string, error) {
	var ret string

	err := c.con.Call("results.get_testjob_results_yaml", id, &ret)

	return ret, err
}

func (c LavaConnection) LavaResults(id int) (LavaResult, error) {
	var ret LavaResult
	yamlStr, err := c.LavaResultsAsYAML(id)

	err = yaml.Unmarshal([]byte(yamlStr), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c LavaConnection) LavaResultsAsJSON(id int) (string, error) {
	results, err := c.LavaResults(id)

	d, err := json.Marshal(&results)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
