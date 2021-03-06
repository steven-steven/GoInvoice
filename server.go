package main

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/steven-steven/GoInvoice/invoice"
	"github.com/steven-steven/GoInvoice/item"
	"github.com/steven-steven/GoInvoice/customer"
)

func newHTTPServer(ctx context.Context, endpoints combinedEndpoint) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/invoice").Handler(httptransport.NewServer(
		endpoints.PostInvoiceEndpoint,
		invoice.DecodePostInvoiceRequest,
		invoice.EncodeResponse,
	))

	r.Methods("GET").Path("/invoice/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}").Handler(httptransport.NewServer(
		endpoints.GetInvoiceEndpoint,
		invoice.DecodeGetInvoiceRequest,
		invoice.EncodeResponse,
	))

	r.Methods("PUT").Path("/invoice/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}").Handler(httptransport.NewServer(
		endpoints.PutInvoiceEndpoint,
		invoice.DecodePutInvoiceRequest,
		invoice.EncodeResponse,
	))
	
	r.Methods("DELETE").Path("/invoice/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}").Handler(httptransport.NewServer(
		endpoints.DeleteInvoiceEndpoint,
		invoice.DecodeDeleteInvoiceRequest,
		invoice.EncodeResponse,
	))

	r.Methods("GET").Path("/allInvoice").Handler(httptransport.NewServer(
		endpoints.GetAllInvoiceEndpoint,
		invoice.DecodeGetAllInvoiceRequest,
		invoice.EncodeResponse,
	))
    
	// ITEM

	r.Methods("POST").Path("/item").Handler(httptransport.NewServer(
		endpoints.PostItemEndpoint,
		item.DecodePostItemRequest,
		item.EncodeResponse,
	))
	
	r.Methods("DELETE").Path("/item/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}").Handler(httptransport.NewServer(
		endpoints.DeleteItemEndpoint,
		item.DecodeDeleteItemRequest,
		item.EncodeResponse,
	))

	r.Methods("GET").Path("/allItem").Handler(httptransport.NewServer(
		endpoints.GetAllItemEndpoint,
		item.DecodeGetAllItemRequest,
		item.EncodeResponse,
	))
    
	// CUSTOMER

	r.Methods("POST").Path("/customer").Handler(httptransport.NewServer(
		endpoints.PostCustomerEndpoint,
		customer.DecodePostCustomerRequest,
		customer.EncodeResponse,
	))

	r.Methods("GET").Path("/customer/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}").Handler(httptransport.NewServer(
		endpoints.GetCustomerEndpoint,
		customer.DecodeGetCustomerRequest,
		customer.EncodeResponse,
	))

	r.Methods("PUT").Path("/customer/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}").Handler(httptransport.NewServer(
		endpoints.PutCustomerEndpoint,
		customer.DecodePutCustomerRequest,
		customer.EncodeResponse,
	))

	r.Methods("DELETE").Path("/customer/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}").Handler(httptransport.NewServer(
		endpoints.DeleteCustomerEndpoint,
		customer.DecodeDeleteCustomerRequest,
		customer.EncodeResponse,
	))

	r.Methods("GET").Path("/allCustomer").Handler(httptransport.NewServer(
		endpoints.GetAllCustomerEndpoint,
		customer.DecodeGetAllCustomerRequest,
		customer.EncodeResponse,
	))
	
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

