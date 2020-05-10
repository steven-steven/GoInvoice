package invoice

import (
    "context"
    "encoding/json"
    "net/http"
)

type postInvoiceRequest struct{
	Invoice
}

type postInvoiceResponse struct {
    Invoice		Invoice_db
}

type getInvoiceRequest struct {
    ID			int
}

type getInvoiceResponse struct {
    Invoice		Invoice_db
}

type putInvoiceRequest struct{
	ID			int
	Invoice
}

type putInvoiceResponse struct {
    Invoice		Invoice_db
}

type getAllInvoiceRequest struct{}

type getAllInvoiceResponse struct {
    Invoices 	map[int]Invoice_db
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

func decodeGetAllInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
