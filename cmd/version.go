package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var cloudBackupCliVersion = ""

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cloudbackup-cli",
	Long:  `Print the version number of cloudbackup-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cloudbackup-cli %s", cloudBackupCliVersion)
	},
}
