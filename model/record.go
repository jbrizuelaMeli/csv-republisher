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
