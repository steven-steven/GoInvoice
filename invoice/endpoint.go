package invoice

import (
    "context"

    "github.com/go-kit/kit/endpoint"
)

// Endpoints are exposed
type Endpoints struct {
    PostInvoiceEndpoint     endpoint.Endpoint
    GetInvoiceEndpoint   	endpoint.Endpoint
	PutInvoiceEndpoint 		endpoint.Endpoint
	GetAllInvoiceEndpoint 	endpoint.Endpoint
}

//Make Endpoints
func MakePostInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(postInvoiceRequest)
        inv, err := srv.PostInvoice(ctx, req.Invoice)
        if err != nil {
            return postInvoiceResponse{}, err
        }
        return postInvoiceResponse{Invoice: inv}, nil
    }
}

func MakeGetInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(getInvoiceRequest)
        inv, err := srv.GetInvoice(ctx, req.ID)
        if err != nil {
            return getInvoiceResponse{}, err
        }
        return getInvoiceResponse{Invoice: inv}, nil
    }
}

func MakePutInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(putInvoiceRequest)
        res, err := srv.PutInvoice(ctx, req.ID, req.Invoice)
        if err != nil {
            return putInvoiceResponse{}, err
        }
        return putInvoiceResponse{Invoice: res}, nil
    }
}

func MakeGetAllInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        _ = request.(getAllInvoiceRequest)
        res, err := srv.GetAllInvoice(ctx)
        if err != nil {
            return getAllInvoiceResponse{}, err
        }
        return getAllInvoiceResponse{Invoices: res}, nil
    }
}