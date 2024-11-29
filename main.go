package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running myapp...")
	},
}

func main() {
	rootCmd.AddCommand(ClientCmd())
	rootCmd.AddCommand(ServerCmd())
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
