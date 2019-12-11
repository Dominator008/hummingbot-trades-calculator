package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	trades    = flag.String("trades", "trades.csv", "Trades CSV")
	fee       = flag.Float64("fee", 0.000405, "Fee rate")
	baseName  = flag.String("base", "ONE", "Name of the base currency")
	quoteName = flag.String("quote", "USDT", "Name of the quote currency")
)

func main() {
	flag.Parse()
	csvFile, err := os.Open(*trades)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var deltaBase float64
	var deltaQuote float64
	var finalPrice float64
	var totalBuyQuantity float64
	var totalSellQuantity float64
	var totalBuyVolume float64
	var totalSellVolume float64
	for i, line := range lines {
		if i == 0 {
			continue
		}
		var isBuy bool
		if line[5] == "buy" {
			isBuy = true
		}
		price, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		if i == len(lines)-1 {
			finalPrice = price
		}
		quantity, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		if isBuy {
			deltaBase += quantity
			deltaQuote -= (1.0 + *fee) * price * quantity
			totalBuyVolume += price * quantity
			totalBuyQuantity += quantity
		} else {
			deltaBase -= quantity
			deltaQuote += (1.0 - *fee) * price * quantity
			totalSellVolume += price * quantity
			totalSellQuantity += quantity
		}
	}
	effectiveDeltaQuote := deltaBase*finalPrice + deltaQuote
	averageBuyPrice := totalBuyVolume / totalBuyQuantity
	averageSellPrice := totalSellVolume / totalSellQuantity
	fmt.Printf("Buy quantity: %.3f %s Sell quantity: %.3f %s\n", totalBuyQuantity, *baseName, totalSellQuantity, *baseName)
	fmt.Printf("Buy volume: %.3f %s Sell volume: %.3f %s\n", totalBuyVolume, *quoteName, totalSellVolume, *quoteName)
	fmt.Printf("Average buy price: %.8f %s/%s\n", averageBuyPrice, *baseName, *quoteName)
	fmt.Printf("Average sell price: %.8f %s/%s\n", averageSellPrice, *baseName, *quoteName)
	fmt.Printf("Delta base: %.3f %s\n", deltaBase, *baseName)
	fmt.Printf("Delta quote: %.3f %s\n", deltaQuote, *quoteName)
	fmt.Printf("Effective delta quote: %.3f %s\n", effectiveDeltaQuote, *quoteName)
}
