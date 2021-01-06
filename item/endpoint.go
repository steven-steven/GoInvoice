package item

import (
    "context"
	"log"
    "github.com/go-kit/kit/endpoint"
)

// Endpoints are exposed
type Endpoints struct {
	PostItemEndpoint    endpoint.Endpoint
	DeleteItemEndpoint 	endpoint.Endpoint
	GetAllItemEndpoint 	endpoint.Endpoint
}

//Make Endpoints
func MakePostItemEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(postItemRequest)
        item, err := srv.PostItem(ctx, req.Item)
        if err != nil {
            return postItemResponse{}, err
		}
		log.Println("post new")
        return postItemResponse{Item: item}, nil
    }
}

func MakeDeleteItemEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(deleteItemRequest)
        res, err := srv.DeleteItem(ctx, req.ID)
        if err != nil {
            return deleteItemResponse{}, err
        }
		log.Printf("delete %v", req.ID)
        return deleteItemResponse{Success: res}, nil
    }
}

func MakeGetAllItemEndpoint(srv Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        _ = request.(getAllItemRequest)
        res, err := srv.GetAllItem(ctx)
        if err != nil {
            return getAllItemResponse{}, err
		}
		log.Println("getAll")
        return getAllItemResponse{Items: res}, nil
    }
}