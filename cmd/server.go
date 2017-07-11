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
	"github.com/spf13/viper"
	"os"
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
		fmt.Println("api client", apiClient)
		fmt.Println("server id", viper.GetInt("server-id"))
	},
}

var serverUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update a server in Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server update called")
		fmt.Println("api client", apiClient)
		fmt.Println("server id", viper.GetInt("server-id"))
	},
}

var serverDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete or remove a server in Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server delete called")
		fmt.Println("api client", apiClient)
		fmt.Println("server id", viper.GetInt("server-id"))
	},
}

var serverInfo = &cobra.Command{
	Use:   "info",
	Short: "Get information for a server in Binlogic CloudBackup",
	Run: func(cmd *cobra.Command, args []string) {
		if server, err := apiClient.GetServer(viper.GetInt("server-id")); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting server.\n%s", err)
		} else {
			if viper.GetBool("json") {
				fmt.Println(server.JSONString())
			} else {
				fmt.Println(server)
			}

		}
	},
}

var serverInstall = &cobra.Command{
	Use:     "install",
	Short:   "Get install link for a server in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("server install called")
		fmt.Println("api client", apiClient)
		fmt.Println("server id", viper.GetInt("server-id"))

		if install, err := apiClient.GetServerInstall(viper.GetInt("server-id")); err != nil {
			return err
		} else {
			if viper.GetBool("dry-run") {
				fmt.Println(string(install))
			} else {
				//TODO actually execute the install script
			}
		}

		return nil
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

	addPersistentInt(serverCmd, "server-id", 0, "Server ID")

	addFlagBool(serverInfo, "json", false, "Output info in JSON format")
	addFlagBool(serverInstall, "dry-run", false, "Output install script instead of executing it")
	// serverInstall.MarkFlagRequired("dry-run")

	// viper.BindPFlag("accesskey", RootCmd.PersistentFlags().Lookup("accesskey"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
