/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bjornm82/trading-obs/internal/strategies"
	"github.com/bjornm82/trading-obs/pkg/broker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	dateFromFlag = "date_from"
	minutesFlag  = "minutes"
)

// ordererCmd represents the orderer command
var ordererCmd = &cobra.Command{
	Use:   "orderer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("orderer called")

		if err := errors.New("api_url should be set"); viper.GetString(apiUrlFlag) == "" {
			log.Fatalln(err)
		}
		if err := errors.New("stream_api_url should be set"); viper.GetString(streamApiUrlFlag) == "" {
			log.Fatalln(err)
		}
		if err := errors.New("token should be set"); viper.GetString(tokenFlag) == "" {
			log.Fatalln(err)
		}
		if err := errors.New("account should be set"); viper.GetString(accountFlag) == "" {
			log.Fatalln(err)
		}
		if err := errors.New("application should be set"); viper.GetString(applicationFlag) == "" {
			log.Fatalln(err)
		}

		cl := broker.NewClient(
			viper.GetString(accountFlag),
			viper.GetString(tokenFlag),
			false,
		)

		e := strategies.Env{
			Broker: broker.New(cl),
		}

		// TODO: 4 o'clock is 6 o'clock on the API and 9 at home? Strange, needs some attention
		// then, err := time.Parse(layout, "2022-07-06T04:00:00")
		// if err != nil {
		// 	return err
		// }

		var from time.Time
		var to time.Time

		if viper.GetString(dateFromFlag) == "" {
			// Set Time to NEW YORK
			to = time.Now().Add(-time.Hour * 4)
			from = to.Add(time.Minute * time.Duration(-viper.GetInt(minutesFlag)))
		} else {
			var layout = "2006-01-02T15:04:05"
			then, err := time.Parse(layout, viper.GetString(dateFromFlag))
			if err != nil {
				fmt.Errorf("unable to parse date %s with error given: %s", viper.GetString(dateFromFlag), err)
				return
			}
			to = then.Add(-time.Hour * 6)
			from = to.Add(time.Minute * time.Duration(-viper.GetInt(minutesFlag)))
		}

		err := e.Run(from, to)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("orderer finished")
	},
}

func init() {
	rootCmd.AddCommand(ordererCmd)

	ordererCmd.Flags().String(apiUrlFlag, "", "Host for server")
	ordererCmd.Flags().String(streamApiUrlFlag, "", "Port for server")
	ordererCmd.Flags().String(tokenFlag, "", "Host for server")
	ordererCmd.Flags().String(accountFlag, "", "Port for server")
	ordererCmd.Flags().String(applicationFlag, "", "Host for server")
	ordererCmd.Flags().Bool(backtestFlag, true, "Port for server")

	ordererCmd.Flags().String(dateFromFlag, "", "Date from where range would start with format [2022-07-06T04:00:00]")
	ordererCmd.Flags().Int(minutesFlag, 60, "Amount of minutes the range would apply")

	viper.BindPFlags(ordererCmd.Flags())

	viper.BindEnv(apiUrlFlag, "OANDA_API_URL")
	viper.BindEnv(streamApiUrlFlag, "OANDA_STREAM_API_URL")
	viper.BindEnv(tokenFlag, "OANDA_TOKEN")
	viper.BindEnv(accountFlag, "OANDA_ACCOUNT")
	viper.BindEnv(applicationFlag, "OANDA_APPLICATION")
	viper.BindEnv(backtestFlag, "OANDA_BACKTEST")

	viper.BindEnv(dateFromFlag, "DATE_FROM")
	viper.BindEnv(minutesFlag, "MINUTES")
}
