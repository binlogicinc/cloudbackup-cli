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
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
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

func (c *Client) GetBackupKeys() (body []byte, err error) {
	resp, err := c.httpClient.SignedGet(c.host+"/backups/keys", defaultHeaders)

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
