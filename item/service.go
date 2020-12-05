package item

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
	"github.com/jeremyschlatter/firebase/db"
	"github.com/steven-steven/GoInvoice/utils"
)

type Service interface {
	PostItem(ctx context.Context, inv Item) (Item_db, error)
	GetItem(ctx context.Context, id string) (Item_db, error)
	DeleteItem(ctx context.Context, id string) (bool, error)
	GetAllItem(ctx context.Context) (map[string]Item_db, error)
}

// DB Model for Item
type Item_db struct {
	Item
	ID			string	`json:"id"`
	CreatedAt	string	`json:"createdAt"`
}

type Item struct {
	Name		string	`json:"name"`
	Description	string	`json:"defaultDesc"`
	Rate		*uint64	`json:"rate"`
}

type itemService struct{
	dbClient	db.Client
}

func NewService(dbClient db.Client) Service {
    return itemService{dbClient}
}

var idGenerator = utils.GenerateUUID

// --- Services ---
var (
	ApiError = errors.New("API Error")
	mux_incrementId sync.Mutex
)

func (srv itemService) PostItem(ctx context.Context, item Item) (Item_db, error) {
	dbClient := srv.dbClient

	now := time.Now()
	itemId := idGenerator()
	acc := Item_db{
		Item: item,
		ID: itemId, // generate uuid
		CreatedAt: now.Format("02/01/2006"),
	}

	if err := dbClient.NewRef("invoice/items/"+itemId).Set(ctx, acc); err != nil {
		log.Println(err)
		return Item_db{}, ApiError
	}

	return acc, nil
}

func (srv itemService) GetItem(ctx context.Context, id string) (Item_db, error) {
    dbClient := srv.dbClient
	
	var res Item_db
	if err := dbClient.NewRef("invoice/documents/"+id).Get(ctx, &res); (err != nil || res.ID == "") {
		return Item_db{}, ApiError
	}
    return res, nil
}

func (srv itemService) DeleteItem(ctx context.Context, id string) (bool, error) {
	dbClient := srv.dbClient

	if err := dbClient.NewRef("invoice/items/"+id).Delete(ctx); err != nil {
		log.Println(err)
		return false, ApiError
	}

    return true, nil
}

func (srv itemService) GetAllItem(ctx context.Context) (map[string]Item_db, error) {
	dbClient := srv.dbClient

	var result map[string]Item_db
	if err := dbClient.NewRef("invoice/items/").Get(ctx, &result); err != nil {
		log.Println(err)
		return map[string]Item_db{}, ApiError
	}
	if (result == nil){
		return map[string]Item_db{}, nil
	}

    return result, nil
}
