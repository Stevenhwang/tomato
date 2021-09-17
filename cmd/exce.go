package cmd

import (
	"tomato/modules"
	"tomato/utils"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "调用和执行 module",

	Run: func(cmd *cobra.Command, args []string) {
		mode, _ := cmd.Flags().GetString("mode")
		hosts, _ := cmd.Flags().GetString("hosts")
		module, _ := cmd.Flags().GetString("module")
		if !utils.FindValInSlice([]string{"ssh", "server"}, mode) {
			cmd.Usage()
			utils.PrintRed("指定的 mode 参数错误")
		}
		groups := utils.ListHosts(hosts)
		moduleList := utils.ListModules()
		if !utils.FindValInSlice(moduleList, module) {
			cmd.Usage()
			utils.PrintGreen("可用模块列表: %v\n", moduleList)
			utils.PrintRed("指定的模块 %s 不存在", module)
		}
		switch module {
		case "ping":
			modules.ExecPing(mode, groups)
		}
	},
}

func init() {
	execCmd.Flags().StringP("mode", "M", "ssh", "连接模式[ssh|server]")
	execCmd.Flags().StringP("hosts", "H", "all", "主机[all|主机组|主机]")
	execCmd.Flags().StringP("module", "m", "shell", "模块")
	RootCmd.AddCommand(execCmd)
}
