package item

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

type postItemRequest struct{
	Item
}

type postItemResponse struct {
    Item		Item_db
}

type getItemRequest struct {
    ID			string
}

type getItemResponse struct {
    Item		Item_db
}

type putItemRequest struct{
	ID			string
	Item
}

type putItemResponse struct {
    Item		Item_db
}

type deleteItemRequest struct{
	ID			string
}

type deleteItemResponse struct {
    Success		bool
}

type getAllItemRequest struct{}

type getAllItemResponse struct {
    Items 	map[string]Item_db
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
    idParam := vars["id"]
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
