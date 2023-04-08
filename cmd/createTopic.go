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
	"github.com/spf13/cobra"
	"github.com/ugol/jr/producers/kafka"
	"log"
)

// createTopicCmd represents the createTopic command
var createTopicCmd = &cobra.Command{
	Use:   "createTopic [topic]",
	Short: "simple command to create a Kafka Topic",
	Long:  "simple command to create a Kafka Topic",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kafkaConfig, _ := cmd.Flags().GetString("kafkaConfig")

		kManager := &kafka.KafkaManager{}
		err := kManager.Initialize(kafkaConfig)
		if err != nil {
			log.Fatal(err)
		}
		partitions, _ := cmd.Flags().GetInt("partitions")
		replica, _ := cmd.Flags().GetInt("replica")
		kManager.CreateTopicFull(args[0], partitions, replica)

	},
}

func init() {
	rootCmd.AddCommand(createTopicCmd)
	createTopicCmd.Flags().IntP("partitions", "p", 6, "Number of partitions")
	createTopicCmd.Flags().IntP("replica", "r", 3, "Replica Factor")
	createTopicCmd.Flags().StringP("kafkaConfig", "F", "./kafka/config.properties", "Kafka configuration")

}
