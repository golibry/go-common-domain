package main

import (
	"encoding/json"
	"fmt"
	"time"

	t "github.com/golibry/go-common-domain/domain/temporal"
)

func main() {
	dateRange, _ := t.NewDateRangeFromString("2026-05-01", "2026-05-03")
	fromDB := t.ReconstituteDateRange(dateRange.StartDate(), dateRange.EndDate())

	jsonData, _ := json.Marshal(dateRange)
	checkDate := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)

	fmt.Println(dateRange.String())
	fmt.Println(dateRange.Contains(checkDate))
	fmt.Println(dateRange.Days())
	fmt.Println(dateRange.Equals(fromDB))
	fmt.Println(string(jsonData))
}
