/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"

	config "github.com/karthikic/techblogs/configs"
	"github.com/karthikic/techblogs/core"
	"github.com/spf13/cobra"
)

var sourceFile string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		entryPoint(args)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func entryPoint(args []string) {
	db := core.SetupDB()
	stop := make(chan bool)

	go cleanUp(stop)

	fmt.Println("Starting the scheduler")
	core.NewScheduler(db, stop)
}

func cleanUp(stop chan bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	for name, _ := range config.GetSources() {
		fmt.Println("Issuing kill to ", name, " worker")
		stop <- true
	}

}
