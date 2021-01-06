package customer

import (
  "context"
	"log"
  "github.com/go-kit/kit/endpoint"
)

// Endpoints are exposed
type Endpoints struct {
	PostCustomerEndpoint    endpoint.Endpoint
	GetCustomerEndpoint   	endpoint.Endpoint
	PutCustomerEndpoint 		endpoint.Endpoint
	DeleteCustomerEndpoint 	endpoint.Endpoint
	GetAllCustomerEndpoint 	endpoint.Endpoint
}

//Make Endpoints
func MakePostCustomerEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postCustomerRequest)
		customer, err := srv.PostCustomer(ctx, req.Customer)
		if err != nil {
				return postCustomerResponse{}, err
		}
		log.Println("post new")
		return postCustomerResponse{Customer: customer}, nil
	}
}

func MakeGetCustomerEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getCustomerRequest)
		customer, err := srv.GetCustomer(ctx, req.ID)
		if err != nil {
			return getCustomerResponse{}, err
		}
		log.Printf("get %v", req.ID)
		return getCustomerResponse{Customer: customer}, nil
	}
}

func MakePutCustomerEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(putCustomerRequest)
		res, err := srv.PutCustomer(ctx, req.ID, req.Customer)
		if err != nil {
			return putCustomerResponse{}, err
		}
		log.Printf("put %v", req.ID)
		return putCustomerResponse{Customer: res}, nil
	}
}

func MakeDeleteCustomerEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteCustomerRequest)
		res, err := srv.DeleteCustomer(ctx, req.ID)
		if err != nil {
				return deleteCustomerResponse{}, err
		}
		log.Printf("delete %v", req.ID)
		return deleteCustomerResponse{Success: res}, nil
	}
}

func MakeGetAllCustomerEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(getAllCustomerRequest)
		res, err := srv.GetAllCustomer(ctx)
		if err != nil {
			return getAllCustomerResponse{}, err
		}
		log.Println("getAll")
		return getAllCustomerResponse{Customers: res}, nil
  }
}