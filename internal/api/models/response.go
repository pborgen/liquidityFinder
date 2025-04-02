package models

import "github.com/pborgen/liquidityFinder/internal/types"

type Pagination struct {
	CurrentPage  int `json:"currentPage"`
	TotalPages   int `json:"totalPages"`
	TotalResults int `json:"totalResults"`
}


type PairsResponse struct {
	Success    bool       			`json:"success"`
	Data      []types.ModelPair     `json:"data"`
	Pagination Pagination 			`json:"pagination"`
}

type TokenAmountsResponse struct {
	Success    bool       				   `json:"success"`
	Data      []types.ModelTokenAmount     `json:"data"`
	Pagination Pagination 				   `json:"pagination"`
}


type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}
