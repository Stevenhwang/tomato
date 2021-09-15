package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "开启 tomato agent 服务",

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("开启 tomato agent 服务")
	},
}

func init() {
	RootCmd.AddCommand(agentCmd)
}
