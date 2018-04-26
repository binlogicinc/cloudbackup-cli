// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"
	"github.com/binlogicinc/cloudbackup-cli/api"
	"github.com/spf13/cobra"
	"os"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Create, update, remove and get information for schedules in Binlogic CloudBackup",
}

var scheduleNew = &cobra.Command{
	Use:     "new",
	Short:   "Add new schedule to Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := getStringFlag(cmd, "name")
		hours := getStringFlag(cmd, "hours")
		days := getStringFlag(cmd, "days")

		sType := getStringFlag(cmd, "schedule-type")
		scheduleType, err := api.ParseScheduleType(sType)

		if err != nil {
			return err
		}

		schedule, err := getAPIClient().CreateSchedule(name, scheduleType, hours, days)

		if err != nil {
			return err
		}

		printVerbose("Schedule created successfully")

		if getBoolFlag(cmd, "json") {
			fmt.Println(schedule.JSONString())
		} else {
			fmt.Println(schedule)
		}

		return nil
	},
}

var scheduleUpdate = &cobra.Command{
	Use:     "update",
	Short:   "Updates a schedule in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		scheduleID := getIntFlag(cmd, "schedule-id")

		if scheduleID == 0 {
			return fmt.Errorf("Schedule ID cannot be zero")
		}

		schedule, err := getAPIClient().GetSchedule(scheduleID)

		if err != nil {
			return err
		}

		if flag := cmd.Flag("name"); flag != nil {
			schedule.Name = flag.Value.String()
		}

		if flag := cmd.Flag("days"); flag != nil {
			schedule.ScheduleDays = getStringFlag(cmd, "days")
		}

		if flag := cmd.Flag("hours"); flag != nil {
			schedule.ScheduleHours = getStringFlag(cmd, "hours")
		}

		if flag := cmd.Flag("schedule-type"); flag != nil {
			newScheduleType, err := api.ParseScheduleType(flag.Value.String())

			if err != nil {
				return err
			}

			schedule.ScheduleType = newScheduleType
		}

		if err := getAPIClient().UpdateSchedule(schedule); err != nil {
			return err
		}

		if getBoolFlag(cmd, "json") {
			fmt.Println(schedule.JSONString())
		} else {
			fmt.Println(schedule)
		}

		return nil
	},
}

var scheduleDelete = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a schedule in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		scheduleID := getIntFlag(cmd, "schedule-id")

		if scheduleID == 0 {
			return fmt.Errorf("Schedule ID cannot be zero")
		}

		if err := getAPIClient().DeleteSchedule(scheduleID); err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		}

		return nil
	},
}

var scheduleInfo = &cobra.Command{
	Use:     "info",
	Short:   "Get information for a schedule in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		scheduleID := getIntFlag(cmd, "schedule-id")

		if scheduleID == 0 {
			return fmt.Errorf("Schedule ID cannot be zero")
		}

		if schedule, err := getAPIClient().GetSchedule(scheduleID); err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		} else {
			if getBoolFlag(cmd, "json") {
				fmt.Println(schedule.JSONString())
			} else {
				fmt.Println(schedule)
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(scheduleCmd)
	scheduleCmd.AddCommand(scheduleNew)
	scheduleCmd.AddCommand(scheduleUpdate)
	scheduleCmd.AddCommand(scheduleDelete)
	scheduleCmd.AddCommand(scheduleInfo)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	scheduleInfo.Flags().Int("schedule-id", 0, "Schedule ID")
	scheduleInfo.MarkFlagRequired("schedule-id")
	scheduleInfo.Flags().Bool("json", false, "Output info in JSON format")

	scheduleDelete.Flags().Int("schedule-id", 0, "Schedule ID")
	scheduleDelete.MarkFlagRequired("schedule-id")

	addCreateScheduleFlags(scheduleNew)

	addCreateScheduleFlags(scheduleUpdate)
	scheduleUpdate.Flags().Int("schedule-id", 0, "Schedule ID")
	scheduleUpdate.MarkFlagRequired("schedule-id")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addCreateScheduleFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("json", false, "Output info in JSON format")

	cmd.Flags().String("name", "", "The schedule name to show in the control panel")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("schedule-type", "", "The schedule type (ondemand, hourly, daily, weekly or monthly)")
	cmd.MarkFlagRequired("schedule-type")

	cmd.Flags().String("hours", "", "For hourly schedule, every how many hours it should run."+
		" For others, at what time of the day will run (00:00 format)")

	cmd.Flags().String("days", "", "For weekly schedule, which days of the week to run (comma "+
		"separated, starting with 0 being Sunday). For monthly, which day of the month to run.")
}
