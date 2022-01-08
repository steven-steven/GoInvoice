package invoice

import (
	"context"
	"errors"
	"log"
	"math"
	"sync"
	"time"

	"github.com/jeremyschlatter/firebase/db"
	"github.com/steven-steven/GoInvoice/utils"
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
	ID        string  `json:"id"`
	CreatedAt string  `json:"createdAt"`
	Total     *uint64 `json:"total"`
	SubTotal  *uint64 `json:"subtotal"`
}

type Invoice struct {
	InvoiceNo          string  `json:"invoice_no"`
	CustomerId         string  `json:"customerId"`
	CatatanInvoice     string  `json:"catatanInvoice"`
	CatatanKwitansi    string  `json:"catatanKwitansi"`
	KeteranganKwitansi string  `json:"keteranganKwitansi"`
	Date               string  `json:"date"`
	Items              []Item  `json:"items"`
	Tax                *uint64 `json:"tax,omitempty"`
}

type Item struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rate        *uint64 `json:"rate"`
	IsMetric    bool    `json:"isMetric"`
	Unit        string  `json:"unit"`
	Quantity    *uint64 `json:"quantity"`
	Amount      *uint64 `json:"amount"`
}

type invoiceService struct {
	dbClient db.Client
}

func NewService(dbClient db.Client) Service {
	return invoiceService{dbClient}
}

var idGenerator = utils.GenerateUUID

// --- Services ---
var (
	ApiError        = errors.New("API Error")
	mux_incrementId sync.Mutex
)

func (srv invoiceService) PostInvoice(ctx context.Context, inv Invoice) (Invoice_db, error) {
	dbClient := srv.dbClient

	// Calculate total: sum all items + tax
	var total uint64
	var subtotal uint64
	var quantity float64
	for i, item := range inv.Items {
		quantity = float64(*item.Quantity)
		if item.IsMetric {
			quantity /= 1000
		}
		calculatedAmount := (quantity * float64(*item.Rate))
		intAmount := uint64(math.Round(calculatedAmount))
		inv.Items[i].Amount = &intAmount
		subtotal += intAmount
	}
	total = subtotal + uint64(math.Round((float64(*inv.Tax)/100)*float64(subtotal)))

	now := time.Now()
	invoiceId := idGenerator()

	acc := Invoice_db{
		Invoice:   inv,
		ID:        invoiceId,
		CreatedAt: now.Format("02/01/2006"),
		Total:     &total,
		SubTotal:  &subtotal,
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
	if err := dbClient.NewRef("invoice/documents/"+id).Get(ctx, &res); err != nil || res.ID == "" {
		log.Println(err)
		return Invoice_db{}, ApiError
	}
	if res.Items == nil {
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
	var subtotal uint64
	var quantity float64
	for i, item := range inv.Items {
		quantity = float64(*item.Quantity)
		if item.IsMetric {
			quantity /= 1000
		}
		calculatedAmount := (quantity * float64(*item.Rate))
		intAmount := uint64(math.Round(calculatedAmount))
		inv.Items[i].Amount = &intAmount
		subtotal += intAmount
	}
	total = subtotal + uint64(math.Round((float64(*inv.Tax)/100)*float64(subtotal)))

	newRecord := Invoice_db{
		Invoice:   inv,
		ID:        id,
		CreatedAt: now.Format("02/01/2006"),
		Total:     &total,
		SubTotal:  &subtotal,
	}

	if err := dbClient.NewRef("invoice/documents/"+id).Set(ctx, newRecord); err != nil {
		log.Println(err)
		return Invoice_db{}, ApiError
	}
	return newRecord, nil
}

func (srv invoiceService) DeleteInvoice(ctx context.Context, id string) (bool, error) {
	dbClient := srv.dbClient

	if err := dbClient.NewRef("invoice/documents/" + id).Delete(ctx); err != nil {
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
	if result == nil {
		return map[string]Invoice_db{}, nil
	}
	for k, inv := range result {
		if inv.Items == nil {
			//https://github.com/golang/go/issues/3117
			var tmp = result[k]
			tmp.Items = make([]Item, 0)
			result[k] = tmp
		}
	}

	return result, nil
}
