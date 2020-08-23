package invoice

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

type postInvoiceRequest struct{
	Invoice
}

type postInvoiceResponse struct {
    Invoice		Invoice_db
}

type getInvoiceRequest struct {
    ID			string
}

type getInvoiceResponse struct {
    Invoice		Invoice_db
}

type putInvoiceRequest struct{
	ID			string
	Invoice
}

type putInvoiceResponse struct {
    Invoice		Invoice_db
}

type deleteInvoiceRequest struct{
	ID			string
}

type deleteInvoiceResponse struct {
    Success		bool
}

type getAllInvoiceRequest struct{}

type getAllInvoiceResponse struct {
    Invoices 	map[string]Invoice_db
}

// Models to JSON
func DecodePostInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req postInvoiceRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        return nil, err
    }
    return req, nil
}

func DecodeGetInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req getInvoiceRequest
    vars := mux.Vars(r)
    idParam := vars["id"]
    req.ID = idParam
    return req, nil
}

func DecodePutInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req putInvoiceRequest
    vars := mux.Vars(r)
    idParam := vars["id"]
    err := json.NewDecoder(r.Body).Decode(&req)
    req.ID = idParam
    if err != nil {
        return nil, err
    }
    return req, nil
}

func DecodeDeleteInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req deleteInvoiceRequest
    vars := mux.Vars(r)
    idParam := vars["id"]
    req.ID = idParam
    return req, nil
}

func DecodeGetAllInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req getAllInvoiceRequest
    return req, nil
}

// JSON to Models
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response)
}
