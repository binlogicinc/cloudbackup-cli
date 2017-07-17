// Copyright © 2017  Alejandro Bednarik <alejandro@binlogic.net>
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
	"os"

	"errors"
	"github.com/binlogicinc/cloudbackup-cli/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string
var apiClient *api.Client

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cloudbackup-cli",
	Short: "command-line tool to interact with Binlogic CloudBackup [ https://www.binlogic.io/ ]",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initAPIClient)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file in toml format (default is $HOME/.cloudbackup-cli.toml)")

	addPersistentString("access-key", "", "API access key", RootCmd)
	addPersistentString("secret-key", "", "API secret key", RootCmd)
	addPersistentString("host", "", "Your host/domain of cloudbackup panel", RootCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cloudbackup-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cloudbackup-cli")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("BL")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initAPIClient() {
	var err error

	accessKey := viper.GetString("access-key")
	secretKey := viper.GetString("secret-key")
	host := viper.GetString("host")

	apiClient, err = api.NewAPIClient(host, accessKey, secretKey)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getBoolFlag(cmd *cobra.Command, name string) bool {
	b, err := cmd.Flags().GetBool(name)

	if err != nil {
		return false
	}

	return b
}

func getIntFlag(cmd *cobra.Command, name string) int {
	i, err := cmd.Flags().GetInt(name)

	if err != nil {
		return 0
	}

	return i
}

func getStringFlag(cmd *cobra.Command, name string) string {
	s, err := cmd.Flags().GetString(name)

	if err != nil {
		return ""
	}

	return s
}

func addPersistentInt(name string, value int, usage string, cmd *cobra.Command) {
	cmd.PersistentFlags().Int(name, value, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func addPersistentBool(name string, value bool, usage string, cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(name, value, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func addPersistentString(name string, value string, usage string, cmd *cobra.Command) {
	cmd.PersistentFlags().String(name, value, usage)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func addFlagInt(name string, value int, usage string, cmd *cobra.Command) {
	cmd.Flags().Int(name, value, usage)
	viper.BindPFlag(name, cmd.Flags().Lookup(name))
}

func addFlagBool(name string, value bool, usage string, cmd *cobra.Command) {
	cmd.Flags().Bool(name, value, usage)
	viper.BindPFlag(name, cmd.Flags().Lookup(name))
}

func addFlagString(name string, value string, usage string, cmd *cobra.Command) {
	cmd.Flags().String(name, value, usage)
	viper.BindPFlag(name, cmd.Flags().Lookup(name))
}

func checkRequiredFlags(cmd *cobra.Command, args []string) error {
	requiredError := false
	flagName := ""

	check := func(flag *pflag.Flag) {
		requiredAnnotation := flag.Annotations[cobra.BashCompOneRequiredFlag]
		if len(requiredAnnotation) == 0 {
			return
		}

		flagRequired := requiredAnnotation[0] == "true"

		if flagRequired && !flag.Changed {
			requiredError = true
			flagName = flag.Name
		}
	}

	cmd.Flags().VisitAll(check)
	cmd.PersistentFlags().VisitAll(check)

	if requiredError {
		return errors.New("Required flag `" + flagName + "` has not been set")
	}

	return nil
}

func checkRoot() error {
	if os.Getuid() != 0 || os.Getgid() != 0 {
		return fmt.Errorf("You need root privileges to execute this command")
	}

	return nil
}
