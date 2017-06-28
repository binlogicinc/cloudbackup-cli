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

import ()

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

func (s *Server) create(httpClient *signedHTTPClient) {

}
