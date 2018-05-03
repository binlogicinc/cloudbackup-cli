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

type Server struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	DbType   databaseType `json:"dbTypeId"`
	Readonly bool         `json:"readonly"`
	DbHost   string       `json:"dbHost"`
	DbPort   string       `json:"dbPort"`
	DbUser   string       `json:"dbUser"`
	DbPass   string       `json:"dbPass"`
}

type databaseType int

const (
	DB_MYSQL    databaseType = 1
	DB_MONGO    databaseType = 2
	DB_POSTGRES databaseType = 3
)

func (d databaseType) String() string {
	switch d {
	case DB_MONGO:
		return "MongoDB"

	case DB_MYSQL:
		return "MySQL"

	case DB_POSTGRES:
		return "PostgreSQL"
	}

	return "Unknown"
}

func ParseDatabaseType(s string) (databaseType, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "mongodb", "mongo":
		return DB_MONGO, nil

	case "mysql", "mariadb", "percona_server":
		return DB_MYSQL, nil

	case "postgresql", "postgre_sql", "postgre", "postgres":
		return DB_POSTGRES, nil
	}

	return 0, fmt.Errorf("Database type %s not recognized", s)
}

func (s Server) String() string {
	return fmt.Sprintf("ID: %d\nName: %s\nDB Type: %s\nReadonly: %t\nDB Host: %s\n"+
		"DB Port: %s\nDB User: %s\nDB Pass: %s\n", s.ID, s.Name, s.DbType.String(),
		s.Readonly, s.DbHost, s.DbPort, s.DbUser, s.DbPass)
}

func (s Server) JSONString() string {
	bs, _ := json.Marshal(s)

	return string(bs)
}

func (c *Client) CreateServer(name string, dbType databaseType, readonly bool,
	dbHost, dbPort, dbUser, dbPass string) (server Server, err error) {

	server = Server{
		0, name, dbType, readonly, dbHost, dbPort, dbUser, dbPass,
	}

	if name == "" {
		return server, fmt.Errorf("Server name cannot be empty")
	}

	if dbHost == "" {
		return server, fmt.Errorf("Database host cannot be empty")
	}

	if dbPort == "" {
		return server, fmt.Errorf("Database port cannot be empty")
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

	if s.Name == "" {
		return fmt.Errorf("Server name cannot be empty")
	}

	if s.DbHost == "" {
		return fmt.Errorf("Database host cannot be empty")
	}

	if s.DbPort == "" {
		return fmt.Errorf("Database port cannot be empty")
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
