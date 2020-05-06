package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/steven-steven/GoInvoice/invoice"
)

func main() {
    var (
        httpAddr = flag.String("http", ":8080", "http listen address")
    )
    flag.Parse()
	ctx := context.Background()
	
	// INVOICE SERVICE
    srv := invoice.NewService()
    errChan := make(chan error)

    go func() {	// cntrl-C
        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        errChan <- fmt.Errorf("%s", <-c)
    }()

    // mapping endpoints
    endpoints := invoice.Endpoints{
		PostInvoiceEndpoint:	invoice.MakePostInvoiceEndpoint(srv),
		GetInvoiceEndpoint   	invoice.MakeGetInvoiceEndpoint(srv),
		PutInvoiceEndpoint 		invoice.MakePutInvoiceEndpoint(srv),
		GetAllInvoiceEndpoint 	invoice.MakeGetAllInvoiceEndpoint(srv),
	}

    // HTTP transport
    go func() {
        log.Println("app listening on port:", *httpAddr)
        handler := invoice.NewHTTPServer(ctx, endpoints)
        errChan <- http.ListenAndServe(*httpAddr, handler)
    }()

    log.Fatalln(<-errChan)
}
