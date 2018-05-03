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

type Schedule struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	ScheduleType  scheduleType `json:"scheduleType"`
	ScheduleHours string       `json:"scheduleHours"`
	ScheduleDays  string       `json:"scheduleDays"`
}

type scheduleType int

const (
	SCHEDULE_ON_DEMAND scheduleType = 1
	SCHEDULE_HOURLY    scheduleType = 2
	SCHEDULE_DAILY     scheduleType = 3
	SCHEDULE_WEEKLY    scheduleType = 4
	SCHEDULE_MONTHLY   scheduleType = 5
)

func (d scheduleType) String() string {
	switch d {
	case SCHEDULE_ON_DEMAND:
		return "On Demand"
	case SCHEDULE_HOURLY:
		return "Hourly"
	case SCHEDULE_DAILY:
		return "Daily"
	case SCHEDULE_WEEKLY:
		return "Weekly"
	case SCHEDULE_MONTHLY:
		return "Monthly"
	}

	return "Unknown"
}

func ParseScheduleType(s string) (scheduleType, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "ondemand":
		return SCHEDULE_ON_DEMAND, nil
	case "hourly":
		return SCHEDULE_HOURLY, nil
	case "daily":
		return SCHEDULE_DAILY, nil
	case "weekly":
		return SCHEDULE_WEEKLY, nil
	case "monthly":
		return SCHEDULE_MONTHLY, nil
	}

	return 0, fmt.Errorf("Schedule type %s not recognized", s)
}

func (s Schedule) String() string {
	daysHours := ""

	if s.ScheduleDays != "" {
		daysHours += "Days: " + s.ScheduleDays

		if s.ScheduleHours != "" {
			daysHours += "\n"
		}
	}

	if s.ScheduleHours != "" {
		daysHours += "Hours: " + s.ScheduleHours
	}

	return fmt.Sprintf("ID: %d\nName: %s\nSchedule Type: %s\n%s",
		s.ID, s.Name, s.ScheduleType, daysHours)
}

func (s Schedule) JSONString() string {
	bs, _ := json.Marshal(s)

	return string(bs)
}

func (c *Client) GetSchedule(id int) (schedule Schedule, err error) {
	resp, err := c.httpClient.SignedGet(c.host+"/schedules/"+strconv.Itoa(id), defaultHeaders)

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

		return schedule, fmt.Errorf("Schedule returned HTTP %d but there is no error "+
			"in response '%s' (this should not happen!)", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &schedule)

	return
}

func (c *Client) DeleteSchedule(id int) error {
	resp, err := c.httpClient.SignedDelete(c.host+"/schedules/"+strconv.Itoa(id), defaultHeaders)

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

func (c *Client) CreateSchedule(name string, scheduleType scheduleType, hours string,
	days string) (schedule Schedule, err error) {

	schedule = Schedule{
		0, name, scheduleType, hours, days,
	}

	val, err := c.httpClient.postJSON(c.host+"/schedules", schedule)

	if err != nil {
		err = wrap("while doing client post", err)
		return
	}

	if id, ok := val["id"]; ok {
		if intID, ok2 := id.(float64); ok2 { //json marshalling converts ints to floats
			schedule.ID = int(intID)
		}
	}

	if schedule.ID <= 0 {
		err = fmt.Errorf("Missing ID from schedule response %v", val)
	}

	return
}

func (c *Client) UpdateSchedule(s Schedule) error {
	if s.ID <= 0 {
		return fmt.Errorf("Invalid ID %d for schedule", s.ID)
	}

	if s.Name == "" {
		return fmt.Errorf("Schedule name cannot be empty")
	}

	_, err := c.httpClient.postJSON(c.host+"/schedules/"+strconv.Itoa(s.ID), s)

	if err != nil {
		return wrap("while doing client post", err)
	}

	return nil
}
