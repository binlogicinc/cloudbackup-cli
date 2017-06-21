package main

import (
	"fmt"
	"github.com/binlogicinc/cloudbackup-cli/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
