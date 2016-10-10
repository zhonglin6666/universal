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
	"net/http"
	"os"
	"runtime"

	"github.com/ngaut/log"
	"github.com/spf13/cobra"

	"github.com/zhonglin6666/universal/api"
	"github.com/zhonglin6666/universal/config"
	"github.com/zhonglin6666/universal/util"
)

var mainCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "A test demo",
	Long:  "test demo",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mainCmd args: %v, len: %v\n", args, len(args))
		if len(args) == 1 {
			cmd.Help()
			return
		}

		file, err := cmd.Flags().GetString("config-file")
		if err != nil {
			fmt.Printf("mainCmd Run get config file err: %v\n", err)
			return
		}

		startMain(file)
	},
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		fmt.Printf("%s", http.ListenAndServe(":9999", nil))
	}()

	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func startMain(file string) {
	fmt.Printf("startMain file: %v\n", file)
	if !util.IsFileExists(file) {
		panic("the config not exist, panic return")
	}

	cfg, err := config.LoadConfig(file)
	if err != nil {
		panic(fmt.Sprintf("load config failure, err:%v", err))
	}

	fmt.Println(cfg.GetString("Default", "NodeHost"))

	log.SetLevelByString(cfg.GetString("Log", "Level"))
	log.SetOutputByName(cfg.GetString("Log", "Path"))
	log.SetRotateByDay()

	startAPIServer(cfg.GetInt("Default", "APIPort"))
}

func startAPIServer(port int) error {
	log.Infof("main start api server, listen on port: %v", port)
	hr := api.HandlerRouter(api.MasterRouter)
	http.Handle("/", hr)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func init() {
	mainCmd.Flags().StringP("config-file", "f", "/etc/smart/master.cfg", "config file")
}
