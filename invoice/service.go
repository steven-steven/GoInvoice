package invoice

import (
	"context"
	"github.com/steven-steven/GoInvoice/config"
)

type Service interface {
    PostInvoice(ctx context.Context, inv Invoice) (Invoice, error)
	GetInvoice(ctx context.Context, id string) (Invoice, error)
	PutInvoice(ctx context.Context, inv Invoice) (Invoice, error)
	GetAllInvoice(ctx context.Context) ([]Invoice, error)
}

// DB Model for Invoice
type Invoice struct {
	Client		string	`json:"client"`
	Date      	string  `json:"date"`
	Items 		[]Item 	`json:"items"`
	Tax			uint32	`json:"tax,omitempty"`
	CreatedAt	string	`json:"createdAt"`
}

type Item struct {
	Name     	string	`json:"name"`
	Rate      	uint32	`json:"rate"`
	Quantity 	int		`json:"quantity"`
	Amount		uint32	`json:"amount"`
}

type invoiceService struct{}

func NewService() Service {
    return invoiceService{}
}

// --- Services ---
var (
	ApiError = errors.New("API Error")
)

func (invoiceService) PostInvoice(ctx context.Context, inv Invoice) (Invoice, error) {
	dbClient, err := getDBInstance(ctx)
	if err != nil {
		return Invoice{}, ApiError
	}
	
	//get invoice id
	idRef := dbClient.NewRef("invoice/nextId")
	var id int
	if err := idRef.Get(ctx, &id); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	
	if err := idRef.Set(ctx, ++id); err != nil {
		log.Fatal(err)
		return ApiError
	}
	fmt.Printf("id %v",id)

	acc := Invoice{
		Client: "PT A",
		Date: "24/03/2019",
		Items: []Item{
			Item{
				Name: "Paku",
				Rate: "10000",
				Quantity: "3",
				Amount: "30000",
			},
		},
		Tax: "5000",
		CreatedAt: "30/05/2019",
	}

	if err := dbClient.NewRef("invoice/"+strconv.Itoa(id)).Set(ctx, acc); err != nil {
		log.Fatal(err)
		return ApiError
	}
	fmt.Println("HI")

	return "ok", nil
}

// Get will return today's date
func (invoiceService) GetInvoice(ctx context.Context, id string) (Invoice, error) {
    now := time.Now()
    return now.Format("02/01/2006"), nil
}

// Validate will check if the date today's date
func (invoiceService) PutInvoice(ctx context.Context, inv Invoice) (Invoice, error) {
    _, err := time.Parse("02/01/2006", date)
    if err != nil {
        return false, err
    }
    return true, nil
}

// Validate will check if the date today's date
func (invoiceService) GetAllInvoice(ctx context.Context) ([]Invoice, error) {
    _, err := time.Parse("02/01/2006", date)
    if err != nil {
        return false, err
    }
    return true, nil
}