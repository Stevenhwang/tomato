package cmd

import (
	"log"
	"tomato/master"

	"github.com/spf13/cobra"
)

var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "开启 tomato master 服务",

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("开启 tomato master 服务")
		master.Start()
	},
}

func init() {
	RootCmd.AddCommand(masterCmd)
}
