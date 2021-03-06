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
	"github.com/spf13/cobra"
)

var backupsCmd = &cobra.Command{
	Use:   "backup",
	Short: "Get all successfull backups information",
}

var backupListKeys = &cobra.Command{
	Use:     "keys",
	Short:   "Prints all your backup encryption keys in JSON format",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		keys, err := getAPIClient().GetBackupKeys()

		if err != nil {
			return err
		}

		fmt.Println(string(keys))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(backupsCmd)
	backupsCmd.AddCommand(backupListKeys)

}
