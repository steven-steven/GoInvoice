package invoice

import (
    "context"
    "encoding/json"
    "net/http"
)

type postInvoiceRequest struct{
	Client		string	`json:"client"`
	Date      	string  `json:"date"`
	Items 		[]Item 	`json:"items"`
	Tax			uint32	`json:"tax,omitempty"`
}

type postInvoiceResponse struct {
	ID			int
    Invoice		Invoice
}

type getInvoiceRequest struct {
    ID			int
}

type getInvoiceResponse struct {
    ID			int
    Invoice		Invoice
}

type putInvoiceRequest struct{
	Client		string	`json:"client"`
	Date      	string  `json:"date"`
	Items 		[]Item 	`json:"items"`
	Tax			uint32	`json:"tax,omitempty"`
}

type putInvoiceResponse struct {
    ID			int
    Invoice		Invoice
}

type getAllInvoiceRequest struct{
	Invoice		Invoice
}

type getAllInvoiceResponse struct {
    Invoices 	map[int]Invoice
}

// Models to JSON
func decodePostInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req postInvoiceRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        return nil, err
    }
    return req, nil
}

func decodeGetInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req getInvoiceRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        return nil, err
    }
    return req, nil
}

func decodePutInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req putInvoiceRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        return nil, err
    }
    return req, nil
}

func decodeGetAllInvoiceResponse(ctx context.Context, r *http.Request) (interface{}, error) {
    var req getAllInvoiceRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        return nil, err
    }
    return req, nil
}

// JSON to Models
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response)
}
