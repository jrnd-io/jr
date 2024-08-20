// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jrnd-io/jr/pkg/configuration"
	"github.com/jrnd-io/jr/pkg/constants"
	"github.com/jrnd-io/jr/pkg/functions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var logLevel = constants.DEFAULT_LOG_LEVEL

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
	rootCmd.PersistentFlags().StringVar(&constants.JR_SYSTEM_DIR, "system_dir", "", "JR system dir")
	rootCmd.PersistentFlags().StringVar(&constants.JR_USER_DIR, "user_dir", "", "JR user dir")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log_level", constants.DEFAULT_LOG_LEVEL, "HR Log Level")
}

func initConfig() {

	// setting zerolog level
	zlogLevel, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		zlogLevel = zerolog.PanicLevel
	}
	zerolog.SetGlobalLevel(zlogLevel)

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
	if constants.JR_SYSTEM_DIR == "" {
		constants.JR_SYSTEM_DIR = constants.SYSTEM_DIR
	}
	viper.AddConfigPath(constants.JR_SYSTEM_DIR)

	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Str("file", viper.ConfigFileUsed()).Msg("JR configuration")
	} else {
		log.Error().Err(err).Msg("JR configuration not found")
	}
	err = viper.UnmarshalKey("global", &configuration.GlobalCfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal global configuration")
	}

	//err = viper.UnmarshalKey("emitters", &emitters)
	err = viper.UnmarshalKey("emitters", &emitters2)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal emitter configuration")
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
