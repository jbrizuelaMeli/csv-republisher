package model

import "time"

type NumericID struct {
	ID int64 `json:"id"`
}

type StringID struct {
	ID string `json:"id"`
}

type MultiRequestNumericIDs struct {
	IDs []NumericID `json:"ids"`
}

type MultiResponseNumericIDs struct {
	IDs    []NumericID `json:"ids"`
	Errors []NumericID `json:"errors"`
}

type ItemRequestMap struct {
	Item map[int64][]string
}

type ItemResponseMap struct {
	Success map[int64][]string
	Errors  map[int64][]string
}

type Report struct {
	ID             string
	TotalLines     int
	ProcessedLines int
	ErrorsQuantity int
	DateLastRun    time.Time
}
