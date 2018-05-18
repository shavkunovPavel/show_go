package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ethCsv struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

func PricesEth(sdate string) (prices []*ethCsv, err error) {
	dt_start := uxdate(sdate)
	csvFile, err := os.Open("./static/csv/eth.csv")
	if err != nil {
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		return
	}

	for ix, each := range csvData {
		if ix > 0 {
			pr := ethCsv{}
			tm := uxdate(each[1])
			if tm.After(dt_start) {
				fl, _ := strconv.ParseFloat(each[2], 32)
				dd := uxdate(each[1])
				pr.Date = fmt.Sprintf("%d-%d-%d", dd.Year(), dd.Month(), dd.Day())
				pr.Amount = fl
				prices = append(prices, &pr)
			}
		}
	}
	return
}

func uxdate(ux string) (tm time.Time) {
	i, err := strconv.ParseInt(ux, 10, 64)
	if err != nil {
		panic(err)
	}
	tm = time.Unix(i, 0)
	return
}
