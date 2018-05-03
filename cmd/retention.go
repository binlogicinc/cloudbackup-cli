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

var retentionCmd = &cobra.Command{
	Use:   "retention",
	Short: "Create, update, remove and get information for retention policies in Binlogic CloudBackup",
}

var retentionNew = &cobra.Command{
	Use:     "new",
	Short:   "Add new retention policy to Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := getStringFlag(cmd, "name")
		retType := getStringFlag(cmd, "retention-type")
		count := getIntFlag(cmd, "count")

		retentionType, err := api.ParseRetentionType(retType)

		if err != nil {
			return err
		}

		retention, err := getAPIClient().CreateRetention(name, retentionType, count)

		if err != nil {
			return err
		}

		printVerbose("Retention created successfully")

		if getBoolFlag(cmd, "json") {
			fmt.Println(retention.JSONString())
		} else {
			fmt.Println(retention)
		}

		return nil
	},
}

var retentionUpdate = &cobra.Command{
	Use:     "update",
	Short:   "Updates a retention policy in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		retentionID := getIntFlag(cmd, "retention-id")

		if retentionID == 0 {
			return fmt.Errorf("Retention ID cannot be zero")
		}

		retention, err := getAPIClient().GetRetention(retentionID)

		if err != nil {
			return err
		}

		if flag := cmd.Flag("name"); flag != nil {
			retention.Name = flag.Value.String()
		}

		if flag := cmd.Flag("count"); flag != nil {
			retention.Count = getIntFlag(cmd, "count")
		}

		if flag := cmd.Flag("retention-type"); flag != nil {
			newRetentionType, err := api.ParseRetentionType(flag.Value.String())

			if err != nil {
				return err
			}

			retention.RetentionType = newRetentionType
		}

		if err := getAPIClient().UpdateRetention(retention); err != nil {
			return err
		}

		if getBoolFlag(cmd, "json") {
			fmt.Println(retention.JSONString())
		} else {
			fmt.Println(retention)
		}

		return nil
	},
}

var retentionDelete = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a retention policy in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		retentionID := getIntFlag(cmd, "retention-id")

		if retentionID == 0 {
			return fmt.Errorf("Retention ID cannot be zero")
		}

		if err := getAPIClient().DeleteRetention(retentionID); err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		}

		return nil
	},
}

var retentionInfo = &cobra.Command{
	Use:     "info",
	Short:   "Get information for a retention policy in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		retentionID := getIntFlag(cmd, "retention-id")

		if retentionID == 0 {
			return fmt.Errorf("Retention ID cannot be zero")
		}

		if retention, err := getAPIClient().GetRetention(retentionID); err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		} else {
			if getBoolFlag(cmd, "json") {
				fmt.Println(retention.JSONString())
			} else {
				fmt.Println(retention)
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(retentionCmd)
	retentionCmd.AddCommand(retentionNew)
	retentionCmd.AddCommand(retentionUpdate)
	retentionCmd.AddCommand(retentionDelete)
	retentionCmd.AddCommand(retentionInfo)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	retentionInfo.Flags().Int("retention-id", 0, "Retention ID")
	retentionInfo.MarkFlagRequired("retention-id")
	retentionInfo.Flags().Bool("json", false, "Output info in JSON format")

	retentionDelete.Flags().Int("retention-id", 0, "Retention ID")
	retentionDelete.MarkFlagRequired("retention-id")

	addCreateRetentionFlags(retentionNew)

	addCreateRetentionFlags(retentionUpdate)
	retentionUpdate.Flags().Int("retention-id", 0, "Retention ID")
	retentionUpdate.MarkFlagRequired("retention-id")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addCreateRetentionFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("json", false, "Output info in JSON format")

	cmd.Flags().String("name", "", "The retention policy name to show in the control panel")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("retention-type", "", "The retention type (bydays or bycount)")
	cmd.MarkFlagRequired("retention-type")

	cmd.Flags().Int("count", 0, "The amount of backups to retain (either days or backups, depending on the retention type)")
	cmd.MarkFlagRequired("count")
}
