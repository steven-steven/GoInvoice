package main

import (
	"context"
	"github.com/spf13/viper"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/steven-steven/GoInvoice/config"
	"github.com/steven-steven/GoInvoice/invoice"
	"github.com/steven-steven/GoInvoice/item"
	"github.com/steven-steven/GoInvoice/customer"
)

type invoiceEndpoint = invoice.Endpoints
type itemEndpoint = item.Endpoints
type customerEndpoint = customer.Endpoints

	type combinedEndpoint struct {
	*invoiceEndpoint
	*itemEndpoint
	*customerEndpoint
}

func main() {
	viper.BindEnv("port")
	var httpAddr = viper.GetString("port")
	if httpAddr == "" {
			httpAddr = "8080"
	}
    
	ctx := context.Background()
	errChan := make(chan error)

	dbClient, err := config.GetRealDB(ctx)
	if err != nil {
		errChan <- fmt.Errorf("%s", err)
	}
	
	// INVOICE SERVICE
	srv := invoice.NewService(*dbClient)
	// ITEM SERVICE
	srvItem := item.NewService(*dbClient)
	// CUSTOMER SERVICE
	srvCustomer := customer.NewService(*dbClient)

	go func() {	// cntrl-C
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := combinedEndpoint{
		&invoice.Endpoints{
			PostInvoiceEndpoint:	invoice.MakePostInvoiceEndpoint(srv),
			GetInvoiceEndpoint:   	invoice.MakeGetInvoiceEndpoint(srv),
			PutInvoiceEndpoint: 	invoice.MakePutInvoiceEndpoint(srv),
			DeleteInvoiceEndpoint: 	invoice.MakeDeleteInvoiceEndpoint(srv),
			GetAllInvoiceEndpoint: 	invoice.MakeGetAllInvoiceEndpoint(srv),
		},
		&item.Endpoints{
			PostItemEndpoint:	    item.MakePostItemEndpoint(srvItem),
			DeleteItemEndpoint: 	item.MakeDeleteItemEndpoint(srvItem),
			GetAllItemEndpoint: 	item.MakeGetAllItemEndpoint(srvItem),
		},
		&customer.Endpoints{
			PostCustomerEndpoint:	    customer.MakePostCustomerEndpoint(srvCustomer),
			GetCustomerEndpoint:   		customer.MakeGetCustomerEndpoint(srvCustomer),
			PutCustomerEndpoint: 			customer.MakePutCustomerEndpoint(srvCustomer),
			DeleteCustomerEndpoint: 	customer.MakeDeleteCustomerEndpoint(srvCustomer),
			GetAllCustomerEndpoint: 	customer.MakeGetAllCustomerEndpoint(srvCustomer),
		},
	}

	// HTTP transport
	go func() {
		log.Println("app listening on port:", httpAddr)
		handler := newHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(":" + httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
