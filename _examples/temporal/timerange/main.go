package main

import (
	"encoding/json"
	"fmt"

	t "github.com/golibry/go-common-domain/domain/temporal"
)

func main() {
	timeRange, _ := t.NewTimeRangeFromString("2026-05-01T10:00:00Z", "2026-05-01T12:30:00Z")
	fromDB := t.ReconstituteTimeRange(timeRange.StartTime(), timeRange.EndTime())

	jsonData, _ := json.Marshal(timeRange)

	fmt.Println(timeRange.String())
	fmt.Println(timeRange.Duration())
	fmt.Println(timeRange.Equals(fromDB))
	fmt.Println(string(jsonData))
}
