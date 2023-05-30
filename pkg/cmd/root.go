//Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package cmd

import (
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/functions"
	"log"
	"os"
	"time"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "jr",
	Short: "jr, the data random generator",
	Long:  `jr is a data random generator that helps you in generating quality random data for your needs.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddGroup(&cobra.Group{
		ID:    "resource",
		Title: "Resources",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "server",
		Title: "HTTP Server",
	})
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jr/jrconfig.json)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home := constants.DEFAULT_HOMEDIR
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName("jrconfig")
		viper.SetEnvPrefix(constants.DEFAULT_ENV_PREFIX)
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		log.Println("JR configuration loaded from:", viper.ConfigFileUsed())
	}
	err := viper.UnmarshalKey("global", &configuration.GlobalCfg)
	if err != nil {
		log.Println(err)
	}
	seed := configuration.GlobalCfg.Seed
	if seed != -1 {
		functions.Random.Seed(seed)
		uuid.SetRand(functions.Random)
	} else {
		functions.Random.Seed(time.Now().UTC().UnixNano())
		uuid.SetRand(functions.Random)
	}
}
