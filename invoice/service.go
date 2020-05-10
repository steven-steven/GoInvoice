package invoice

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"
	"github.com/jeremyschlatter/firebase/db"
)

type Service interface {
    PostInvoice(ctx context.Context, inv Invoice) (Invoice_db, error)
	GetInvoice(ctx context.Context, id int) (Invoice_db, error)
	PutInvoice(ctx context.Context, id int, inv Invoice) (Invoice_db, error)
	DeleteInvoice(ctx context.Context, id int) (bool, error)
	GetAllInvoice(ctx context.Context) (map[int]Invoice_db, error)
}

// DB Model for Invoice
type Invoice_db struct {
	Invoice
	ID			int		`json:"id"`
	CreatedAt	string	`json:"createdAt"`
}

type Invoice struct {
	Client		string	`json:"client"`
	Date      	string  `json:"date"`
	Items 		[]Item 	`json:"items"`
	Tax			uint32	`json:"tax,omitempty"`
}

type Item struct {
	Name     	string	`json:"name"`
	Rate      	uint32	`json:"rate"`
	Quantity 	int		`json:"quantity"`
	Amount		uint32	`json:"amount"`
}

type invoiceService struct{
	dbClient	db.Client
}

func NewService(dbClient db.Client) Service {
    return invoiceService{dbClient}
}

// --- Services ---
var (
	ApiError = errors.New("API Error")
)

func (srv invoiceService) PostInvoice(ctx context.Context, inv Invoice) (Invoice_db, error) {
	dbClient := srv.dbClient
	
	//get invoice id
	idRef := dbClient.NewRef("invoice/lastId")
	var id int
	if err := idRef.Get(ctx, &id); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	id++
	if err := idRef.Set(ctx, id); err != nil {
		log.Fatal(err)
		return Invoice_db{}, ApiError
	}

	now := time.Now()
	acc := Invoice_db{
		Invoice: inv,
		ID: id,
		CreatedAt: now.Format("02/01/2006"),
	}

	if err := dbClient.NewRef("invoice/documents/"+strconv.Itoa(id)).Set(ctx, acc); err != nil {
		log.Println(err)
		return Invoice_db{}, ApiError
	}

	return acc, nil
}

func (srv invoiceService) GetInvoice(ctx context.Context, id int) (Invoice_db, error) {
    dbClient := srv.dbClient
	
	var res Invoice_db

	if err := dbClient.NewRef("invoice/documents/"+strconv.Itoa(id)).Get(ctx, &res); (err != nil || res.ID == 0) {
		log.Println(err)
		return Invoice_db{}, ApiError
	}
    return res, nil
}

func (srv invoiceService) PutInvoice(ctx context.Context, id int, inv Invoice) (Invoice_db, error) {
    dbClient := srv.dbClient

	//new data
	now := time.Now()
	newRecord := Invoice_db{
		Invoice: inv,
		ID: id,
		CreatedAt: now.Format("02/01/2006"),
	}

	if err := dbClient.NewRef("invoice/documents/"+strconv.Itoa(id)).Set(ctx, newRecord); err != nil {
		log.Println(err)
		return Invoice_db{}, ApiError
	}
    return newRecord, nil
}

func (srv invoiceService) DeleteInvoice(ctx context.Context, id int) (bool, error) {
	dbClient := srv.dbClient

	if err := dbClient.NewRef("invoice/documents/"+strconv.Itoa(id)).Delete(ctx); err != nil {
		log.Println(err)
		return false, ApiError
	}

    return true, nil
}

func (srv invoiceService) GetAllInvoice(ctx context.Context) (map[int]Invoice_db, error) {
	dbClient := srv.dbClient

	var result map[int]Invoice_db
	if err := dbClient.NewRef("invoice/documents/").Get(ctx, &result); err != nil {
		log.Println(err)
		return map[int]Invoice_db{}, ApiError
	}

    return result, nil
}
