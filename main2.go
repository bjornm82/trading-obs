package main

import (
	"fmt"

	"github.com/xtordoir/goanda/api"
	"github.com/xtordoir/goanda/models"
)

func priceProcessor(c chan models.ClientPrice) {
	for {
		data := <-c
		tick := models.ClientPrice2Tick(&data)
		fmt.Println(tick)
	}
}

func tickProcessor(c chan models.Tick) {
	for {
		tick := <-c
		fmt.Println(tick)
	}
}

func heartbeatProcessor(c chan models.PricingHeartbeat) {
	for {
		data := <-c
		fmt.Printf("%s\n", data)
	}
}

func main2() {

	// channels for data
	pchan := make(chan models.ClientPrice)
	tchan := make(chan models.Tick)
	hchan := make(chan models.PricingHeartbeat)

	// start processors for data
	go priceProcessor(pchan)
	go tickProcessor(tchan)
	go heartbeatProcessor(hchan)

	// context to create api
	ctx := api.Context{
		ApiURL:       "https://api-fxpractice.oanda.com",
		StreamApiURL: "https://stream-fxpractice.oanda.com/",
		Token:        "a7f4f34b8f41a1d1d467047517c9b8a0-73ef8ea2c18a28fbb89c0fc21786b9b9",
		Account:      "101-004-8979399-001",
		Application:  "TestStreaming",
	}

	fmt.Printf("%s\n", ctx.ApiURL)
	fmt.Printf("%s\n", ctx.StreamApiURL)
	// fmt.Printf("%s\n", ctx.Token)
	fmt.Printf("%s\n", ctx.Account)
	fmt.Printf("%s\n", ctx.Application)

	api := ctx.CreateAPI()

	if len(ctx.Account) == 0 {
		accounts, err := api.GetAccounts()
		if err == nil && len(accounts.Accounts) > 0 {
			fmt.Printf("Setting Account # in context: %s\n", accounts.Accounts[0].ID)
			ctx.Account = accounts.Accounts[0].ID
			api = ctx.CreateAPI()
		}
	}

	pos, err := api.GetPricing([]string{"EUR_USD"})

	// err := api.GetAccounts()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	fmt.Printf("%s\n", pos)

	streamapi := ctx.CreateStreamAPI()
	//streamapi.PricingStream([]string{"EUR_USD", "BCO_USD", "SPX500_USD", "EUR_JPY"}, pchan, hchan)

	streamapi.TickStream([]string{"EUR_USD", "BCO_USD", "SPX500_USD", "EUR_JPY"}, tchan, hchan)
}
