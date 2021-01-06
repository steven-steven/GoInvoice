package customer

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
  PostCustomer(ctx context.Context, customer Customer) (Customer_db, error)
	GetCustomer(ctx context.Context, id string) (Customer_db, error)
	PutCustomer(ctx context.Context, id string, customer Customer) (Customer_db, error)
	DeleteCustomer(ctx context.Context, id string) (bool, error)
	GetAllCustomer(ctx context.Context) (map[string]Customer_db, error)
}

// DB Model for Customer
type Customer_db struct {
	Customer
	ID			string	`json:"id"`
	CreatedAt	string	`json:"createdAt"`
}

type Customer struct {
	Client		string	`json:"client"`
	ClientAddress *ClientAddress `json:"client_address,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type ClientAddress struct {
	Address		string	`json:"address,omitempty"`
	City		string	`json:"city,omitempty"`
	State      	string  `json:"state,omitempty"`
	Country 	string 	`json:"country,omitempty"`
	PostalCode	string	`json:"postal_code,omitempty"`
}

type customerService struct{
	dbClient	db.Client
}

func NewService(dbClient db.Client) Service {
	return customerService{dbClient}
}

var idGenerator = utils.GenerateUUID

// --- Services ---
var (
	ApiError = errors.New("API Error")
	mux_incrementId sync.Mutex
)

func (srv customerService) PostCustomer(ctx context.Context, customer Customer) (Customer_db, error) {
	dbClient := srv.dbClient

	now := time.Now()
	customerId := idGenerator()

	acc := Customer_db{
		Customer: customer,
		ID: customerId,
		CreatedAt: now.Format("02/01/2006"),
	}

	if err := dbClient.NewRef("invoice/customers/"+customerId).Set(ctx, acc); err != nil {
		log.Println(err)
		return Customer_db{}, ApiError
	}

	return acc, nil
}

func (srv customerService) GetCustomer(ctx context.Context, id string) (Customer_db, error) {
	dbClient := srv.dbClient
	
	var res Customer_db
	if err := dbClient.NewRef("invoice/customers/"+id).Get(ctx, &res); (err != nil || res.ID == "") {
		return Customer_db{}, ApiError
	}
	return res, nil
}

func (srv customerService) PutCustomer(ctx context.Context, id string, customer Customer) (Customer_db, error) {
	dbClient := srv.dbClient

	//new data
	now := time.Now()
	
	newRecord := Customer_db{
		Customer: customer,
		ID: id,
		CreatedAt: now.Format("02/01/2006"),
	}

	if err := dbClient.NewRef("invoice/customers/"+id).Set(ctx, newRecord); err != nil {
		return Customer_db{}, ApiError
	}
	return newRecord, nil
}

func (srv customerService) DeleteCustomer(ctx context.Context, id string) (bool, error) {
	dbClient := srv.dbClient

	if err := dbClient.NewRef("invoice/customers/"+id).Delete(ctx); err != nil {
		log.Println(err)
		return false, ApiError
	}

	return true, nil
}

func (srv customerService) GetAllCustomer(ctx context.Context) (map[string]Customer_db, error) {
	dbClient := srv.dbClient

	var result map[string]Customer_db
	if err := dbClient.NewRef("invoice/customers/").Get(ctx, &result); err != nil {
		log.Println(err)
		return map[string]Customer_db{}, ApiError
	}
	if (result == nil){
		return map[string]Customer_db{}, nil
	}
	return result, nil
}
