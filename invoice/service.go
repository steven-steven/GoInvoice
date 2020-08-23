package invoice

import (
	"context"
	"fmt"
	"errors"
	"log"
	"sync"
	"time"
	"github.com/jeremyschlatter/firebase/db"
)

type Service interface {
    PostInvoice(ctx context.Context, inv Invoice) (Invoice_db, error)
	GetInvoice(ctx context.Context, id string) (Invoice_db, error)
	PutInvoice(ctx context.Context, id string, inv Invoice) (Invoice_db, error)
	DeleteInvoice(ctx context.Context, id string) (bool, error)
	GetAllInvoice(ctx context.Context) (map[string]Invoice_db, error)
}

// DB Model for Invoice
type Invoice_db struct {
	Invoice
	ID			string	`json:"id"`
	CreatedAt	string	`json:"createdAt"`
	Total		*uint64	`json:"total,omitempty"`
}

type Invoice struct {
	Client		string	`json:"client"`
	ClientAddress *ClientAddress `json:"client_address,omitempty"`
	Date      	string  `json:"date"`
	Items 		[]Item 	`json:"items"`
	Tax			*uint64	`json:"tax,omitempty"`
}

type ClientAddress struct {
	Address		string	`json:"address,omitempty"`
	City		string	`json:"city,omitempty"`
	State      	string  `json:"state,omitempty"`
	Country 	string 	`json:"country,omitempty"`
	PostalCode	string	`json:"postal_code,omitempty"`
}

type Item struct {
	Name     	string	`json:"name"`
	Rate      	*uint64	`json:"rate"`
	Quantity 	int		`json:"quantity"`
	Amount		*uint64	`json:"amount"`
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
	mux_incrementId sync.Mutex
)

func (srv invoiceService) PostInvoice(ctx context.Context, inv Invoice) (Invoice_db, error) {
	dbClient := srv.dbClient
	
	t, _ := time.Parse("02/01/2006", inv.Date)
	dateId := t.Format("0601")

	//get invoice id
	mux_incrementId.Lock()
	idRef := dbClient.NewRef("invoice/lastId/"+dateId)
	var id int
	if err := idRef.Get(ctx, &id); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	id++
	if err := idRef.Set(ctx, id); err != nil {
		log.Fatal(err)
		return Invoice_db{}, ApiError
	}
	mux_incrementId.Unlock()

	now := time.Now()

	// Calculate total: sum all items + tax
	var total uint64
	for i, item := range inv.Items {
		// replace amount value
		calculatedQuantity := (uint64(item.Quantity) * (*item.Rate))
		inv.Items[i].Amount = &calculatedQuantity
		total += calculatedQuantity
	}
	total += (*inv.Tax)

	invoiceId := dateId+"-"+fmt.Sprintf("%05d", id)
	acc := Invoice_db{
		Invoice: inv,
		ID: invoiceId,
		CreatedAt: now.Format("02/01/2006"),
		Total: &total,
	}

	if err := dbClient.NewRef("invoice/documents/"+invoiceId).Set(ctx, acc); err != nil {
		log.Println(err)
		return Invoice_db{}, ApiError
	}

	return acc, nil
}

func (srv invoiceService) GetInvoice(ctx context.Context, id string) (Invoice_db, error) {
    dbClient := srv.dbClient
	
	var res Invoice_db
	if err := dbClient.NewRef("invoice/documents/"+id).Get(ctx, &res); (err != nil || res.ID == "") {
		log.Println(err)
		return Invoice_db{}, ApiError
	}
	if(res.Items == nil){
		res.Items = make([]Item, 0)
	}
    return res, nil
}

func (srv invoiceService) PutInvoice(ctx context.Context, id string, inv Invoice) (Invoice_db, error) {
    dbClient := srv.dbClient

	//new data
	now := time.Now()

	// Calculate total: sum all items + tax
	var total uint64
	for i, item := range inv.Items {
		// replace amount value
		calculatedQuantity := (uint64(item.Quantity) * (*item.Rate))
		inv.Items[i].Amount = &calculatedQuantity
		total += calculatedQuantity
	}
	total += (*inv.Tax)

	newRecord := Invoice_db{
		Invoice: inv,
		ID: id,
		CreatedAt: now.Format("02/01/2006"),
		Total: &total,
	}

	if err := dbClient.NewRef("invoice/documents/"+id).Set(ctx, newRecord); err != nil {
		log.Println(err)
		return Invoice_db{}, ApiError
	}
    return newRecord, nil
}

func (srv invoiceService) DeleteInvoice(ctx context.Context, id string) (bool, error) {
	dbClient := srv.dbClient

	if err := dbClient.NewRef("invoice/documents/"+id).Delete(ctx); err != nil {
		log.Println(err)
		return false, ApiError
	}

    return true, nil
}

func (srv invoiceService) GetAllInvoice(ctx context.Context) (map[string]Invoice_db, error) {
	dbClient := srv.dbClient

	var result map[string]Invoice_db
	// https://github.com/golang/go/issues/37711 (nill instead of []. For now need to manually go over and explicitly make [])
	if err := dbClient.NewRef("invoice/documents/").Get(ctx, &result); err != nil {
		log.Println(err)
		return map[string]Invoice_db{}, ApiError
	}
	for k, inv := range result {
        if(inv.Items == nil){
			//https://github.com/golang/go/issues/3117
			var tmp = result[k]
			tmp.Items = make([]Item, 0)
			result[k] = tmp
		}
	}

    return result, nil
}
