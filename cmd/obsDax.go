/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/bjornm82/trading-obs/internal/strategies"
	"github.com/bjornm82/trading-obs/pkg/broker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	apiUrlFlag       = "api_url"
	streamApiUrlFlag = "stream_api_url"
	tokenFlag        = "token"
	accountFlag      = "account"
	applicationFlag  = "application"
	backtestFlag     = "backtest"
)

// obsDaxCmd represents the obsDax command
var obsDaxCmd = &cobra.Command{
	Use:   "obsDax",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("obsDax called")

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

		ok, err := e.ObsPosition()
		if err != nil {
			log.Fatalln(err)
		}

		if !ok {
			log.Fatalln("something went wrong, not sure what")
		} else {
			log.Println("all processed")
		}

		fmt.Println("obsDax finished")
	},
}

func init() {
	rootCmd.AddCommand(obsDaxCmd)

	obsDaxCmd.Flags().String(apiUrlFlag, "", "Host for server")
	obsDaxCmd.Flags().String(streamApiUrlFlag, "", "Port for server")
	obsDaxCmd.Flags().String(tokenFlag, "", "Host for server")
	obsDaxCmd.Flags().String(accountFlag, "", "Port for server")
	obsDaxCmd.Flags().String(applicationFlag, "", "Host for server")
	obsDaxCmd.Flags().Bool(backtestFlag, true, "Port for server")

	viper.BindPFlags(obsDaxCmd.Flags())

	viper.BindEnv(apiUrlFlag, "OANDA_API_URL")
	viper.BindEnv(streamApiUrlFlag, "OANDA_STREAM_API_URL")
	viper.BindEnv(tokenFlag, "OANDA_TOKEN")
	viper.BindEnv(accountFlag, "OANDA_ACCOUNT")
	viper.BindEnv(applicationFlag, "OANDA_APPLICATION")
	viper.BindEnv(backtestFlag, "OANDA_BACKTEST")
}
