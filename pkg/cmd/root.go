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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/functions"
	"log"
	"os"
	"strings"
	"time"
)

var home string

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
	rootCmd.PersistentFlags().StringVar(&home, "config", "", "JR config dir")
}

func initConfig() {

	viper.SetConfigName("jrconfig")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	for _, path := range strings.Split(os.ExpandEnv("$PATH"), ":") {
		viper.AddConfigPath(path)
	}
	viper.SetEnvPrefix(constants.DEFAULT_ENV_PREFIX)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	bindFlags(rootCmd, viper.GetViper())
	viper.AddConfigPath(home)

	if err := viper.ReadInConfig(); err == nil {
		log.Println("JR configuration loaded from:", viper.ConfigFileUsed())
	} else {
		log.Println(err)
	}
	err := viper.UnmarshalKey("global", &configuration.GlobalCfg)
	if err != nil {
		log.Println(err)
	}
	err = viper.UnmarshalKey("emitters", &emitters)
	if err != nil {
		log.Println(err)
	}
	seed := configuration.GlobalCfg.Seed
	if seed != -1 {
		functions.SetSeed(seed)
	} else {
		functions.SetSeed(time.Now().UTC().UnixNano())
	}
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
