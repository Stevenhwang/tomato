package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "调用和执行 modules",

	Run: func(cmd *cobra.Command, args []string) {
		hosts, _ := cmd.Flags().GetString("hosts")
		modules, _ := cmd.Flags().GetString("modules")
		log.Println(hosts)
		log.Println(modules)
		log.Println(args)
	},
}

func init() {
	execCmd.Flags().StringP("hosts", "H", "all", "主机")
	execCmd.Flags().StringP("modules", "m", "shell", "模块名")
	RootCmd.AddCommand(execCmd)
}
