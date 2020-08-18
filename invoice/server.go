package invoice

import (
    "context"
    "net/http"

    httptransport "github.com/go-kit/kit/transport/http"
    "github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
    r := mux.NewRouter()
    r.Use(commonMiddleware)

    r.Methods("POST").Path("/invoice").Handler(httptransport.NewServer(
        endpoints.PostInvoiceEndpoint,
        decodePostInvoiceRequest,
        encodeResponse,
    ))

    r.Methods("GET").Path("/invoice").Handler(httptransport.NewServer(
        endpoints.GetInvoiceEndpoint,
        decodeGetInvoiceRequest,
        encodeResponse,
    ))

    r.Methods("PUT").Path("/invoice").Handler(httptransport.NewServer(
        endpoints.PutInvoiceEndpoint,
        decodePutInvoiceRequest,
        encodeResponse,
	))
	
	r.Methods("DELETE").Path("/invoice/{id:[0-9]+}").Handler(httptransport.NewServer(
        endpoints.DeleteInvoiceEndpoint,
        decodeDeleteInvoiceRequest,
        encodeResponse,
    ))

	r.Methods("GET").Path("/allInvoice").Handler(httptransport.NewServer(
        endpoints.GetAllInvoiceEndpoint,
        decodeGetAllInvoiceRequest,
        encodeResponse,
	))
	
    return r
}

func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

