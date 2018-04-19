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
