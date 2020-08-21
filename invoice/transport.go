package invoice

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "strconv"
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

type deleteInvoiceRequest struct{
	ID			int
}

type deleteInvoiceResponse struct {
    Success		bool
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
    vars := mux.Vars(r)
    idParam, err := strconv.Atoi(vars["id"])
    if err != nil {
        return nil, err
    }
    req.ID = idParam
    return req, nil
}

func decodePutInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req putInvoiceRequest
    vars := mux.Vars(r)
    idParam, err := strconv.Atoi(vars["id"])
    if err != nil {
        return nil, err
    }
    err = json.NewDecoder(r.Body).Decode(&req)
    req.ID = idParam
    if err != nil {
        return nil, err
    }
    return req, nil
}

func decodeDeleteInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req deleteInvoiceRequest
    vars := mux.Vars(r)
    idParam, err := strconv.Atoi(vars["id"])
    if err != nil {
        return nil, err
    }
    req.ID = idParam
    return req, nil
}

func decodeGetAllInvoiceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req getAllInvoiceRequest
    return req, nil
}

// JSON to Models
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response)
}
