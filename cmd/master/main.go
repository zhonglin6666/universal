// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zhonglin6666/universal/imp"
)

var cfgFile = "/home/lin/test.cfg"
var name string
var age int

var GlobalViper = viper.New()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "demo",
	Short: "A test demo",
	Long:  "test demo",

	Run: func(cmd *cobra.Command, args []string) {
		if len(name) == 0 {
			cmd.Help()
			return
		}
		imp.Show(name, age)
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("111111111111 ", GlobalViper.GetString("ContentDir"))
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.demo.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.Flags().StringVarP(&name, "name", "n", "AAAA", "person's name")
	RootCmd.Flags().IntVarP(&age, "age", "a", 0, "person's age")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		GlobalViper.SetConfigFile(cfgFile)
	}

	GlobalViper.SetConfigName("test.cfg") // name of config file (without extension)
	GlobalViper.AddConfigPath("$HOME")    // adding home directory as first search path
	GlobalViper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := GlobalViper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", GlobalViper.ConfigFileUsed())
	}

	GlobalViper.SetDefault("ContentDir", "content")

	fmt.Println("1111111111111111111111111 ", GlobalViper.GetString("ContentDir"))
}
