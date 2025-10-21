package cmd

import (
	"os"

	"github.com/fatih/color"
	tools "github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <target>",
	Short: "Remove a file or folder",
	Long:  `Remove a file or folder.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		color.Cyan("開始刪-除：%s", target)
		defer color.Cyan("結束刪-除：%s", target)

		err := os.RemoveAll(target)
		if err != nil {
			tools.Errorf("%s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
