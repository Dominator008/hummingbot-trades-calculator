# Hummingbot trades calculator

A simple calculator for Hummingbot trades. Prints deltas in the base and quote
currencies, and calculates profitability in terms of the quote currency based on
the final price of the base currency.

The current version assumes a pure market making strategy running on a single
exchange.

Usage:

```
go run calc.go --trades=trades.csv --fee=0.00075 --base=ONE --quote=BNB
```
