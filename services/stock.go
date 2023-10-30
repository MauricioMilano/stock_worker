package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func EvalStock(key string) string {
	stockServiceUrl := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", url.QueryEscape(key))
	log.Println("info : processing", stockServiceUrl)

	response, err := http.Get(stockServiceUrl)
	if err != nil {
		log.Println("error :", err)
		return "Stock service is not available"
	}

	if response.StatusCode == http.StatusOK {
		content, err := csv.NewReader(response.Body).ReadAll()
		if err != nil {
			log.Println("error :", err)
			return "Stock service CSV error"
		}
		stock_name := content[1][0]
		pricing := content[1][6]
		log.Println("content:", content)
		if pricing == "N/D" {
			return fmt.Sprintf("%s quote is not available", strings.ToUpper(stock_name))
		}
		return fmt.Sprintf("%s quote is $%s per share", strings.ToUpper(stock_name), pricing)
	}

	log.Println("error : Status ", response.StatusCode)
	return "Stock service is not available"
}
