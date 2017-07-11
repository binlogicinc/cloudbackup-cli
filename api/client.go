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
	"net/url"
	"path"
	"strconv"
)

var (
	defaultHeaders = map[string]string{
		"Content-type": "application/json",
	}
)

type Client struct {
	host       string
	httpClient *signedHTTPClient
}

func NewAPIClient(host, accessKey, accessSecret string) (*Client, error) {
	if host == "" {
		return nil, fmt.Errorf("API Host cannot be empty")
	}
	if accessKey == "" {
		return nil, fmt.Errorf("API Access key cannot be empty")
	}
	if accessSecret == "" {
		return nil, fmt.Errorf("API Access secret cannot be empty")
	}

	u, err := url.Parse(host)

	if err != nil {
		return nil, err
	}

	u.Scheme = "https"
	u.Path = path.Join(u.Path, "api")

	return &Client{
		host:       u.String(),
		httpClient: NewSignedHTTPClient(accessKey, accessSecret, 10), //default 10 secs timeout
	}, nil
}

func (c *Client) CreateServer(name string, dbType databaseType, readonly bool,
	dbHost, dbPort, dbUser, dbPass string) (server Server, err error) {

	server = Server{
		0, name, dbType, readonly, dbHost, dbPort, dbUser, dbPass,
	}

	val, err := c.httpClient.postJSON(c.host+"/servers", server)

	if err != nil {
		err = wrap("while doing client post", err)
		return
	}

	if id, ok := val["id"]; ok {
		if intID, ok2 := id.(float64); ok2 { //json marshalling converts ints to floats
			server.ID = int(intID)
		}
	}

	if server.ID <= 0 {
		err = fmt.Errorf("Missing ID from server response %v", val)
	}

	return
}

func (c *Client) UpdateServer(s Server) error {
	if s.ID <= 0 {
		return fmt.Errorf("Invalid ID %d for server", s.ID)
	}

	_, err := c.httpClient.postJSON(c.host+"/servers/"+strconv.Itoa(s.ID), s)

	if err != nil {
		return wrap("while doing client post", err)
	}

	return nil
}

func (c *Client) DeleteServer(id int) error {
	resp, err := c.httpClient.SignedDelete(c.host+"/servers/"+strconv.Itoa(id), defaultHeaders)

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

func (c *Client) GetServer(id int) (server Server, err error) {
	resp, err := c.httpClient.SignedGet(c.host+"/servers/"+strconv.Itoa(id), defaultHeaders)

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

		return server, fmt.Errorf("Server returned HTTP %d but there is no error "+
			"in response %s (this should not happen!)", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &server)

	return
}

func (c *Client) GetServerInstall(id int) (body []byte, err error) {
	resp, err := c.httpClient.SignedGet(c.host+"/servers/"+strconv.Itoa(id)+"/install", defaultHeaders)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		err = fmt.Errorf("Server responded HTTP %d: %s", resp.StatusCode, string(body))
	}

	return
}

func (c *Client) String() string {
	return fmt.Sprintf("Host: %s, Access Key: %s, Secret Key %s", c.host,
		c.httpClient.AccessKey, c.httpClient.SecretKey)
}

func wrap(context string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s, %s", err, context)
}
