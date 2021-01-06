package customer

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type postCustomerRequest struct{
	Customer
}

type postCustomerResponse struct {
  Customer		Customer_db
}

type getCustomerRequest struct {
  ID			string
}

type getCustomerResponse struct {
  Customer		Customer_db
}

type putCustomerRequest struct{
	ID			string
	Customer
}

type putCustomerResponse struct {
  Customer		Customer_db
}

type deleteCustomerRequest struct{
	ID			string
}

type deleteCustomerResponse struct {
  Success		bool
}

type getAllCustomerRequest struct{}

type getAllCustomerResponse struct {
  Customers 	map[string]Customer_db
}

// Models to JSON
func DecodePostCustomerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req postCustomerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeGetCustomerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getCustomerRequest
	vars := mux.Vars(r)
	idParam := vars["id"]
	req.ID = idParam
	return req, nil
}

func DecodePutCustomerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req putCustomerRequest
	vars := mux.Vars(r)
	idParam := vars["id"]
	err := json.NewDecoder(r.Body).Decode(&req)
	req.ID = idParam
	if err != nil {
			return nil, err
	}
	return req, nil
}

func DecodeDeleteCustomerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req deleteCustomerRequest
	vars := mux.Vars(r)
	idParam := vars["id"]
	req.ID = idParam
	return req, nil
}

func DecodeGetAllCustomerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getAllCustomerRequest
	return req, nil
}

// JSON to Models
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
