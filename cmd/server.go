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

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Create, update, remove and get information for servers in Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

var serverNew = &cobra.Command{
	Use:   "new",
	Short: "Add new servers to Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server new called")
	},
}

var serverUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update a server in Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server update called")
	},
}

var serverDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete or remove a server in Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server delete called")
	},
}

var serverInfo = &cobra.Command{
	Use:   "info",
	Short: "Get information for a server in Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server info called")
	},
}

var serverInstall = &cobra.Command{
	Use:   "install",
	Short: "Get install link for a server in Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server install called")
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverNew)
	serverCmd.AddCommand(serverUpdate)
	serverCmd.AddCommand(serverDelete)
	serverCmd.AddCommand(serverInfo)
	serverCmd.AddCommand(serverInstall)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
