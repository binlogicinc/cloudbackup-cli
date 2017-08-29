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
	"bytes"
	"fmt"
	"github.com/binlogicinc/cloudbackup-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Create, update, remove and get information for servers in Binlogic CloudBackup",
}

var serverNew = &cobra.Command{
	Use:     "new",
	Short:   "Add new servers to Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := getStringFlag(cmd, "name")
		readonly := getBoolFlag(cmd, "readonly")
		dbType := getStringFlag(cmd, "db-type")
		dbHost := getStringFlag(cmd, "db-host")
		dbPort := getStringFlag(cmd, "db-port")
		dbUser := getStringFlag(cmd, "db-user")
		dbPass := getStringFlag(cmd, "db-pass")

		databaseType, err := api.ParseDatabaseType(dbType)

		if err != nil {
			return err
		}

		server, err := apiClient.CreateServer(name, databaseType, readonly, dbHost, dbPort, dbUser, dbPass)

		if err != nil {
			return err
		}

		fmt.Println("Server created successfully")

		if getBoolFlag(cmd, "json") {
			fmt.Println(server.JSONString())
		} else {
			fmt.Println(server)
		}

		return nil
	},
}

var serverUpdate = &cobra.Command{
	Use:     "update",
	Short:   "Update a server in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverID := getIntFlag(cmd, "server-id")

		if serverID == 0 {
			return fmt.Errorf("Server ID cannot be zero")
		}

		server, err := apiClient.GetServer(serverID)

		if err != nil {
			return err
		}

		if flag := cmd.Flag("name"); flag != nil {
			server.Name = flag.Value.String()
		}

		if flag := cmd.Flag("db-host"); flag != nil {
			server.DbHost = flag.Value.String()
		}

		if flag := cmd.Flag("db-port"); flag != nil {
			server.DbPort = flag.Value.String()
		}

		if flag := cmd.Flag("db-user"); flag != nil {
			server.DbUser = flag.Value.String()
		}

		if flag := cmd.Flag("db-pass"); flag != nil {
			server.DbPass = flag.Value.String()
		}

		if flag := cmd.Flag("readonly"); flag != nil {
			server.Readonly = getBoolFlag(cmd, "readonly")
		}

		if flag := cmd.Flag("db-type"); flag != nil {
			newDbType, err := api.ParseDatabaseType(flag.Value.String())

			if err != nil {
				return err
			}

			if server.DbType != newDbType {
				return fmt.Errorf("Can't change db-type from %s to %s", server.DbType, newDbType)
			}
		}

		if err := apiClient.UpdateServer(server); err != nil {
			return err
		}

		if getBoolFlag(cmd, "json") {
			fmt.Println(server.JSONString())
		} else {
			fmt.Println(server)
		}

		return nil
	},
}

var serverDelete = &cobra.Command{
	Use:     "delete",
	Short:   "Delete or remove a server in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverID := getIntFlag(cmd, "server-id")

		if serverID == 0 {
			return fmt.Errorf("Server ID cannot be zero")
		}

		return apiClient.DeleteServer(serverID)
	},
}

var serverInfo = &cobra.Command{
	Use:     "info",
	Short:   "Get information for a server in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverID := getIntFlag(cmd, "server-id")

		if serverID == 0 {
			return fmt.Errorf("Server ID cannot be zero")
		}

		if server, err := apiClient.GetServer(serverID); err != nil {
			return err
		} else {
			if getBoolFlag(cmd, "json") {
				fmt.Println(server.JSONString())
			} else {
				fmt.Println(server)
			}
		}

		return nil
	},
}

var serverInstall = &cobra.Command{
	Use:     "install",
	Short:   "Install a server in this host or print the install script via stdout",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverID := getIntFlag(cmd, "server-id")

		if serverID == 0 {
			return fmt.Errorf("Server ID cannot be zero")
		}

		if install, err := apiClient.GetServerInstall(serverID); err != nil {
			return err
		} else {
			if viper.GetBool("dry-run") {
				fmt.Println(string(install))
			} else {
				if err := checkRoot(); err != nil {
					return err
				}

				cmd := exec.Command("bash")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = bytes.NewReader(install)

				if err := cmd.Run(); err != nil {
					return err
				}
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

	serverInstall.Flags().Int("server-id", 0, "Server ID")
	serverInstall.Flags().Bool("dry-run", false, "Output install script instead of executing it")

	serverInfo.Flags().Int("server-id", 0, "Server ID")
	serverInfo.Flags().Bool("json", false, "Output info in JSON format")

	serverDelete.Flags().Int("server-id", 0, "Server ID")

	addCreateServerFlags(serverNew)
	serverNew.MarkFlagRequired("name")
	serverNew.MarkFlagRequired("db-type")
	serverNew.MarkFlagRequired("db-port")

	addCreateServerFlags(serverUpdate)
	serverUpdate.Flags().Int("server-id", 0, "Server ID")

	serverDelete.MarkPersistentFlagRequired("server-id")
	serverInstall.MarkPersistentFlagRequired("server-id")
	serverInfo.MarkPersistentFlagRequired("server-id")
	serverUpdate.MarkPersistentFlagRequired("server-id")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addCreateServerFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("json", false, "Output info in JSON format")
	cmd.Flags().String("name", "", "The server name to show in the control panel")
	cmd.Flags().String("db-type", "", "The database type (mysql, mariadb, percona_server, mongodb, postgresql)")
	cmd.Flags().Bool("readonly", false, "If the server is readonly (can be backed up but can't receive restores)")
	cmd.Flags().String("db-host", "localhost", "The host of the database to connect the agent to")
	cmd.Flags().String("db-port", "", "The port of the database to connect the agent to")
	cmd.Flags().String("db-user", "", "The user the agent will use to connect to the database")
	cmd.Flags().String("db-pass", "", "The password the agent will use to connect to the database")
}
