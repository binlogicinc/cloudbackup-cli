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
	case "mongodb":
		return DB_MONGO, nil
	case "mysql":
		return DB_MYSQL, nil
	case "mariadb":
		return DB_MYSQL, nil
	case "percona_server":
		return DB_MYSQL, nil
	case "postgresql":
		return DB_POSTGRES, nil
	}

	return 0, fmt.Errorf("Database type %s not recognized", s)
}
