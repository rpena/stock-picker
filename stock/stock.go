package stock

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"stock-picker/config"
)

// Output data
type StockInfo struct {
	Symbol          string               `json:"Symbol"`
	AvgPrice        float64              `json:"Average Price"`
	TimeSeriesDaily map[string]DailyData `json:"Time Series (Daily)"`
}

// Entire JSON structure
type StockData struct {
	MetaData        MetaData             `json:"Meta Data"`
	TimeSeriesDaily map[string]DailyData `json:"Time Series (Daily)"`
}

// "Meta Data" section
type MetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

// "Time Series (Daily)" section
type DailyData struct {
	Open             string `json:"1. open"`
	High             string `json:"2. high"`
	Low              string `json:"3. low"`
	Close            string `json:"4. close"`
	AdjustedClose    string `json:"5. adjusted close"`
	Volume           string `json:"6. volume"`
	DividendAmount   string `json:"7. dividend amount"`
	SplitCoefficient string `json:"8. split coefficient"`
}

// Return a map of DailyData only containing at most numDays of data
func getNDays(data map[string]DailyData, numDays int) map[string]DailyData {
	if len(data) == 0 {
		return data
	}

	if numDays > len(data) {
		numDays = len(data)
	}

	// Get list of dates
	var dates []string
	for date := range data {
		// TODO: only append if within date range
		dates = append(dates, date)
	}

	// Sort dates since they could be out of order
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	nTimeSeriesDaily := make(map[string]DailyData, numDays)
	for i := range dates {
		date := dates[i]
		nTimeSeriesDaily[date] = data[date]
		if len(nTimeSeriesDaily) == numDays {
			break
		}
	}
	return nTimeSeriesDaily
}

// Return avg closing price of provided map of DailyData
func getAvgPrice(data map[string]DailyData) float64 {
	var sum float64

	if len(data) == 0 {
		return sum
	}

	for _, dailyData := range data {
		value, _ := strconv.ParseFloat(dailyData.AdjustedClose, 64)
		sum += value
	}
	return sum / float64(len(data))
}

func GetStockData() (*StockInfo, error) {
	// apikey, err := config.ReadAPIKey()
	// if err != nil {
	// 	return nil, err
	// }
	apikey := "3XZLSTHJI7CJ44PY"
	//log.Printf(": %s\n", apikey)

	symbol, err := config.GetEnvData(config.SYMBOL_ENV)
	if err != nil {
		return nil, err
	}
	nDays, err := config.GetEnvData(config.NDAYS_ENV)
	if err != nil {
		return nil, err
	}
	numDays, err := strconv.Atoi(nDays)
	if err != nil {
		log.Printf("Unable to convert numDays: %s", err.Error())
		return nil, err
	}

	var stockData StockData

	log.Printf("Getting stock information for stock symbol: %s", symbol)
	resp, err := http.Get("https://www.alphavantage.co/query?apikey=" + apikey + "&function=TIME_SERIES_DAILY_ADJUSTED&symbol=" + symbol)
	if err != nil {
		log.Printf("Error calling GET on external API resource")
		return nil, fmt.Errorf("error getting stock info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Got unexpected status code calling external API resource")
		return nil, fmt.Errorf("external API returned status code: %d", resp.StatusCode)
	}

	// Decode json data into stockData for data handling below
	if err := json.NewDecoder(resp.Body).Decode(&stockData); err != nil {
		log.Printf("Error decoding json from external API resource")
		return nil, fmt.Errorf("error decoding external API response: %w", err)
	}

	nTimeSeriesDaily := getNDays(stockData.TimeSeriesDaily, numDays)
	avgPrice := getAvgPrice(nTimeSeriesDaily)

	stockInfo := &StockInfo{
		Symbol:          symbol,
		AvgPrice:        avgPrice,
		TimeSeriesDaily: nTimeSeriesDaily,
	}
	return stockInfo, nil
}
