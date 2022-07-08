# Trading OBS

## Outbreak strategy

1. TradeSetup
GetDataFromLastHour foreach instrument
GetHighest and Lowest
Calculate the price for Entry, SL and TP
ListenOnSocket
When Price is < > trigger market order with TP and SL
Remove instrument from list

2. TradeManager
Cron every 5 minutes
If a position closes half of the points to take profit, update order with SL to that price

## Conceptual idea:

In the hour when the market opens, the volatility of the price is as high as the range in the hour pre-market.

On the opening minute, place a sell and buy limit-stop order on the top of the range previous hour.
SL is half of the range TP is the full range, resulting with a 2R.

The balance on my account is 20000 and I only want to risk 1% of my account size, resulting in USD 200 per trade.
Price range previous (between 08:00 and 09:00) hour is between 100 and 110.
If time is 09:00 place buy and sell stoplimit order with:

USD 200 per trade with 2R

BUY price: 110 sl: 105 tp: 120
SELL price: 100 sl: 105 tp: 90

The order to be created has a lot size of: 200 / 5 = 40 shares in both situations.

When either the BUY or SELL order is being initiated, the other should be cancelled.

When the closing price of the position is half way, move the SL to break even.