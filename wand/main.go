package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/torniker/codenotary/wand/accounting"
)

func main() {
	var rootCmd = &cobra.Command{Use: "wand"}
	accounting.Register(rootCmd)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
