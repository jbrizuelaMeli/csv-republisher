package model

type NumericIDs struct {
	IDs []int64 `json:"ids"`
}

type NumericID struct {
	ID int64 `json:"id"`
}

type StringIDs struct {
	IDs []string `json:"ids"`
}

type StringID struct {
	ID string `json:"id"`
}

type MultiRequestNumericIDs struct {
	IDs []NumericIDs `json:"ids"`
}

type MultiResponseNumericIDs struct {
	IDs    []NumericIDs `json:"ids"`
	Errors []NumericIDs `json:"errors"`
}
