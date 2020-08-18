package invoice

import (
    "context"
    "testing"
	"time"
	// "fmt"
	// "reflect"
	"github.com/stretchr/testify/assert"
    "github.com/steven-steven/GoInvoice/config"
)

func TestPostInvoice(t *testing.T) {
	srv, ctx := setup()	//new test DB
	
	tests := map[string]struct {
        input  Invoice
        output Invoice_db
        err    error
    }{
        "successful post": {
			input:  Invoice{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"},"24/03/2019",[]Item{Item{"Paku",10000,3,30000},Item{"Dua",5000,2,100}},5000},
            output:	Invoice_db{Invoice{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"},"24/03/2019",[]Item{Item{"Paku",10000,3,30000},Item{"Dua",5000,2,10000}},5000},1,time.Now().Format("02/01/2006"),45000},
           	err:    nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.PostInvoice(ctx, test.input)
			assert.IsType(t, test.err, err)
        	assert.EqualValues(t, test.output, output)
		})
	}
}

func TestGetInvoice(t *testing.T) {
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostInvoice(ctx, Invoice{"PT C",nil,"24/03/2018",[]Item{Item{"Paku",10000,3,30000}},5000})
	srv.PostInvoice(ctx, Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu",10000,3,30000}},6000})
	srv.PostInvoice(ctx, Invoice{"PT A",nil,"24/03/2019",[]Item{Item{"Paku",10000,3,30000}},5000})

	tests := map[string]struct {
        input  int
        output Invoice_db
        err    error
    }{
        "successful get": {
            input:  2,
            output:	Invoice_db{Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu",10000,3,30000}},6000},2,time.Now().Format("02/01/2006"),36000},
           	err:    nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.GetInvoice(ctx, test.input)
			assert.IsType(t, test.err, err)
        	assert.EqualValues(t, test.output, output)
		})
	}
}

func TestPutInvoice(t *testing.T) {
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostInvoice(ctx, Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu",10000,3,30000}},6000})
	srv.PostInvoice(ctx, Invoice{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"},"24/03/2019",[]Item{Item{"Paku",10000,3,30000}},5000})

	tests := map[string]struct {
		input_id	int
        input  		Invoice
        output 		Invoice_db
        err    		error
    }{
        "successful put": {
			input_id:	2,
            input:  	Invoice{"PT C",&ClientAddress{Address:"St",PostalCode:""},"24/03/2019",[]Item{Item{"Paku",10000,3,30000}},5000},
            output:		Invoice_db{Invoice{"PT C",&ClientAddress{Address:"St",PostalCode:""},"24/03/2019",[]Item{Item{"Paku",10000,3,30000}},5000},2,time.Now().Format("02/01/2006"),35000},
           	err:    	nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.PutInvoice(ctx, test.input_id, test.input)
			assert.IsType(t, test.err, err)
			assert.EqualValues(t, test.output, output)
			//check get returns updated
			after, errAfter := srv.GetInvoice(ctx, test.input_id)
			assert.IsType(t, nil, errAfter)
			assert.EqualValues(t, test.output, after)
		})
	}
}

func TestDeleteInvoice(t *testing.T) {
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostInvoice(ctx, Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu",10000,3,30000}},6000})
	srv.PostInvoice(ctx, Invoice{"PT A",nil,"24/03/2019",[]Item{Item{"Paku",10000,3,30000}},5000})

	tests := map[string]struct {
        input  		int
        output 		bool
        err    		error
    }{
        "successful delete": {
			input:	2,
			output: true,
           	err:	nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.DeleteInvoice(ctx, test.input)
			assert.IsType(t, test.err, err)
			assert.EqualValues(t, test.output, output)
			//check get returns updated
			after, errAfter := srv.GetInvoice(ctx, test.input)
			assert.IsType(t, ApiError, errAfter)
			assert.EqualValues(t, Invoice_db{}, after)
		})
	}
}

func TestGetAllInvoice(t *testing.T) {
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostInvoice(ctx, Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu",10000,3,30000}},6000})
	srv.PostInvoice(ctx, Invoice{"PT A",nil,"24/03/2019",[]Item{Item{"Paku",10000,3,30000}},5000})

	tests := map[string]struct {
        output map[int]Invoice_db
        err    error
    }{
        "successful get all": {
			output:	map[int]Invoice_db{
				1: Invoice_db{Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu",10000,3,30000}},6000},1,time.Now().Format("02/01/2006"),36000},
				2: Invoice_db{Invoice{"PT A",nil,"24/03/2019",[]Item{Item{"Paku",10000,3,30000}},5000},2,time.Now().Format("02/01/2006"),35000},
			},
           	err:    nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.GetAllInvoice(ctx)
			assert.IsType(t, test.err, err)
        	assert.EqualValues(t, test.output, output)
		})
	}
}

func setup() (srv Service, ctx context.Context) {
	dbClient := config.GetTestDB(ctx)
    return NewService(dbClient), context.Background()
}