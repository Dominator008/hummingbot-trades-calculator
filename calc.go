package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
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
	var volume float64
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
		volume += quantity
		if isBuy {
			deltaBase += quantity
			deltaQuote -= (1.0 + *fee) * price * quantity
		} else {
			deltaBase -= quantity
			deltaQuote += (1.0 - *fee) * price * quantity
		}
	}
	effectiveDeltaQuote := deltaBase*finalPrice + deltaQuote
	profitPercentage := effectiveDeltaQuote / math.Abs(deltaQuote) * 100
	fmt.Printf("Volume: %.3f %s\n", volume, *baseName)
	fmt.Printf("Delta base: %.3f %s\n", deltaBase, *baseName)
	fmt.Printf("Delta quote: %.3f %s\n", deltaQuote, *quoteName)
	fmt.Printf("Effective delta quote: %.3f %s\n", effectiveDeltaQuote, *quoteName)
	fmt.Printf("Profit percentage %.3f%%\n", profitPercentage)
}
