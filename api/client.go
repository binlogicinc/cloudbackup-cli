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
	"net/url"
	"path"
	"strconv"
)

var (
	defaultHeaders = map[string]string{
		"Content-type": "application/json",
	}
)

type databaseType int

const (
	DB_MYSQL    databaseType = 1
	DB_MONGO    databaseType = 2
	DB_POSTGRES databaseType = 3
)

type client struct {
	host       string
	httpClient *signedHTTPClient
}

func NewAPIClient(host, accessKey, accessSecret string) (*client, error) {
	u, err := url.Parse(host)

	if err != nil {
		return nil, err
	}

	u.Scheme = "https"
	u.Path = path.Join(u.Path, "api")

	return &client{
		host:       u.String(),
		httpClient: NewSignedHTTPClient(accessKey, accessSecret, 10), //default 10 secs timeout
	}, nil
}

func (c *client) CreateServer(name string, dbType databaseType, readonly bool,
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

func (c *client) UpdateServer(s *Server) {
	//POST http://192.168.10.80/api/servers/521/
}

func (c *client) DeleteServer(id int) error {
	// DELETE http://192.168.10.80/api/servers/521/
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

func (c *client) GetServer() {
	//GET http://192.168.10.80/api/servers/521/
}

func wrap(context string, err error) error {
	return fmt.Errorf("%s, %s", err, context)
}
