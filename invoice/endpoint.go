package invoice

import (
    "context"
    "errors"

    "github.com/go-kit/kit/endpoint"
)

// Endpoints are exposed
type Endpoints struct {
    PostInvoiceEndpoint      endpoint.Endpoint
    GetInvoiceEndpoint   	endpoint.Endpoint
	PutInvoiceEndpoint 		endpoint.Endpoint
	GetAllInvoiceEndpoint 	endpoint.Endpoint
}

//Make Endpoints
func MakePostInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req = request.(postInvoiceRequest)
        inv, err := srv.PostInvoice(ctx, req)
        if err != nil {
            return postInvoiceResponse{}, err.Error()
        }
        return postInvoiceResponse{inv}, nil
    }
}

func MakeGetInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req = request.(getInvoiceRequest)
        inv, err := srv.GetInvoice(ctx, req)
        if err != nil {
            return getInvoiceResponse{}, err.Error()
        }
        return getInvoiceResponse{inv}, nil
    }
}

func MakePutInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(putInvoiceRequest)
        res, err := srv.PutInvoice(ctx, req)
        if err != nil {
            return putInvoiceResponse{res}, err.Error()
        }
        return putInvoiceResponse{res}, nil
    }
}

func MakeGetAllInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(getAllInvoiceRequest)
        res, err := srv.GetAllInvoice(ctx, req)
        if err != nil {
            return getAllInvoiceResponse{res}, err.Error()
        }
        return getAllInvoiceResponse{res}, nil
    }
}