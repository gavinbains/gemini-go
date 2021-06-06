# gemini-go
Script to use Gemini API to get tickers and place orders. 

**Note:** USE AT YOUR OWN RISK -- this code is provided as-is and built primarily for my own personal use as a learning opportunity.
## BEFORE YOU RUN
In `main.go` there are 3 variables to change based on your needs.

```go
const buySize = 300.00
const tickerSymbol = "BTCUSD"
const precision = "%.8f"
```

`buySize`: is the dollar amount you wish to purchase

`tickerSymbol`: is the [symbol of the ticker](https://docs.gemini.com/rest-api/?shell#symbols-and-minimums) for the currency you want to purchase.

`precision`: is the precision of the currency i.e. BTC is 8 decimals (`"%.8f"`), ETH is 6 decimals(`"%.6f"`), etc 
## HOW TO RUN
```
go build
./gemini-go
```
