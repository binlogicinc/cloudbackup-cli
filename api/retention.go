// Copyright © 2017  Fermin Silva <fermin@binlogic.net>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Retention struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	RetentionType retentionType `json:"retentionType"`
	Count         int           `json:"count"`
}

type retentionType int

const (
	RETENTION_BY_DAYS  retentionType = 1
	RETENTION_BY_COUNT retentionType = 2
)

func (d retentionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d retentionType) String() string {
	switch d {
	case RETENTION_BY_DAYS:
		return "By Days"

	case RETENTION_BY_COUNT:
		return "By Count"
	}

	return "Unknown"
}

func ParseRetentionType(s string) (retentionType, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "bydays":
		return RETENTION_BY_DAYS, nil

	case "bycount":
		return RETENTION_BY_COUNT, nil
	}

	return 0, fmt.Errorf("Retention type %s not recognized", s)
}

func (r Retention) String() string {
	return fmt.Sprintf("ID: %d\nName: %s\nRetention Type: %s\nCount: %d",
		r.ID, r.Name, r.RetentionType, r.Count)
}

func (r Retention) JSONString() string {
	bs, _ := json.Marshal(r)

	return string(bs)
}

func (c *Client) GetRetention(id int) (retention Retention, err error) {
	resp, err := c.httpClient.SignedGet(c.host+"/retentions/"+strconv.Itoa(id), defaultHeaders)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {

		if err != nil {
			return
		}

		_, err = c.httpClient.isResponseOk(body)

		if err != nil {
			return
		}

		return retention, fmt.Errorf("Retention returned HTTP %d but there is no error "+
			"in response '%s' (this should not happen!)", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &retention)

	return
}

func (c *Client) DeleteRetention(id int) error {
	resp, err := c.httpClient.SignedDelete(c.host+"/retentions/"+strconv.Itoa(id), defaultHeaders)

	if err != nil {
		return err
	}

	body, val, err := c.httpClient.parseResponseJSON(resp)

	if err != nil {
		return err
	}

	_, err = c.httpClient.isJSONResponseOk(body, val)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateRetention(name string, retentionType retentionType,
	count int) (retention Retention, err error) {

	retention = Retention{
		0, name, retentionType, count,
	}

	val, err := c.httpClient.postJSON(c.host+"/retentions", retention)

	if err != nil {
		err = wrap("while doing client post", err)
		return
	}

	if id, ok := val["id"]; ok {
		if intID, ok2 := id.(float64); ok2 { //json marshalling converts ints to floats
			retention.ID = int(intID)
		}
	}

	if retention.ID <= 0 {
		err = fmt.Errorf("Missing ID from retention response %v", val)
	}

	return
}

func (c *Client) UpdateRetention(r Retention) error {
	if r.ID <= 0 {
		return fmt.Errorf("Invalid ID %d for retention", r.ID)
	}

	if r.Name == "" {
		return fmt.Errorf("Retention name cannot be empty")
	}

	if r.Count <= 0 {
		return fmt.Errorf("Retention count cannot be <= 0")
	}

	_, err := c.httpClient.postJSON(c.host+"/retentions/"+strconv.Itoa(r.ID), r)

	if err != nil {
		return wrap("while doing client post", err)
	}

	return nil
}
