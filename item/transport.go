package item

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "strconv"
)

type postItemRequest struct{
	Item
}

type postItemResponse struct {
    Item		Item_db
}

type getItemRequest struct {
    ID			int
}

type getItemResponse struct {
    Item		Item_db
}

type putItemRequest struct{
	ID			int
	Item
}

type putItemResponse struct {
    Item		Item_db
}

type deleteItemRequest struct{
	ID			int
}

type deleteItemResponse struct {
    Success		bool
}

type getAllItemRequest struct{}

type getAllItemResponse struct {
    Items 	map[int]Item_db
}

// Models to JSON
func DecodePostItemRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req postItemRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        return nil, err
    }
    return req, nil
}

func DecodeDeleteItemRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req deleteItemRequest
    vars := mux.Vars(r)
    idParam, err := strconv.Atoi(vars["id"])
    if err != nil {
        return nil, err
    }
    req.ID = idParam
    return req, nil
}

func DecodeGetAllItemRequest(ctx context.Context, r *http.Request) (interface{}, error) {
    var req getAllItemRequest
    return req, nil
}

// JSON to Models
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response)
}
