package invoice

import (
    "context"
	"log"
    "github.com/go-kit/kit/endpoint"
)

// Endpoints are exposed
type Endpoints struct {
    PostInvoiceEndpoint     endpoint.Endpoint
    GetInvoiceEndpoint   	endpoint.Endpoint
	PutInvoiceEndpoint 		endpoint.Endpoint
	DeleteInvoiceEndpoint 	endpoint.Endpoint
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
		log.Println("post new")
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
		log.Printf("get %v", req.ID)
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
		log.Printf("put %v", req.ID)
        return putInvoiceResponse{Invoice: res}, nil
    }
}

func MakeDeleteInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(deleteInvoiceRequest)
        res, err := srv.DeleteInvoice(ctx, req.ID)
        if err != nil {
            return deleteInvoiceResponse{}, err
        }
		log.Printf("delete %v", req.ID)
        return deleteInvoiceResponse{Success: res}, nil
    }
}

func MakeGetAllInvoiceEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        _ = request.(getAllInvoiceRequest)
        res, err := srv.GetAllInvoice(ctx)
        if err != nil {
            return getAllInvoiceResponse{}, err
		}
		log.Println("getAll")
        return getAllInvoiceResponse{Invoices: res}, nil
    }
}