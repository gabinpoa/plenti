package cmd

import (
	"github.com/spf13/cobra"
)

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage themes",
	Long: `Download, enable, update, or remove "themes" that your
site can use to inherit content, layouts, media, and static.`,
}

func init() {
	rootCmd.AddCommand(themeCmd)
}
