package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	font string
)

// wrapCmd represents the wrap command
var wrapCmd = &cobra.Command{
	Use:   "wrap",
	Short: "wrap text",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wrap called")
	},
}

func init() {
	rootCmd.AddCommand(wrapCmd)
	wrapCmd.PersistentFlags().StringVarP(&font, "font", "f", "", "a ttf font to use")
}
