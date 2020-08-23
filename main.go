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
)

type invoiceEndpoint = invoice.Endpoints
type itemEndpoint = item.Endpoints

type combinedEndpoint struct {
    *invoiceEndpoint
    *itemEndpoint
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
	}

    // HTTP transport
    go func() {
        log.Println("app listening on port:", httpAddr)
        handler := newHTTPServer(ctx, endpoints)
        errChan <- http.ListenAndServe(":" + httpAddr, handler)
    }()

    log.Fatalln(<-errChan)
}
