package invoice

import (
    "context"
	"testing"
	"time"
	// "fmt"
	// "reflect"
	"math"
	"github.com/stretchr/testify/assert"
    "github.com/steven-steven/GoInvoice/config"
)

var itemRate1 = uint64(10000)
var itemQuantity1 = 3
var itemAmount1 = itemRate1*uint64(itemQuantity1)
var itemRate2 = uint64(5000)
var itemQuantity2 = 2
var itemAmount2 = itemRate2*uint64(itemQuantity2)
var itemTax1 = uint64(5000)
var itemTax2 = uint64(6000)
var total_1_2_tax1 = (itemAmount1 + itemAmount2) + uint64(math.Round(float64(itemTax1)/100*float64(itemAmount1 + itemAmount2)))
var total_1_2_tax2 = (itemAmount1 + itemAmount2) + uint64(math.Round(float64(itemTax2)/100*float64(itemAmount1 + itemAmount2)))
var total_1_tax1 = itemAmount1 + uint64(math.Round(float64(itemTax1)/100*float64(itemAmount1)))
var total_2_tax2 = itemAmount2 + uint64(math.Round(float64(itemTax2)/100*float64(itemAmount2)))
var total_2_tax1 = itemAmount2 + uint64(math.Round(float64(itemTax1)/100*float64(itemAmount2)))
var total_1_tax2 = itemAmount1 + uint64(math.Round(float64(itemTax2)/100*float64(itemAmount1)))

func TestPostInvoice(t *testing.T) {
	srv, ctx := setup()	//new test DB
	
	tests := map[string]struct {
        input  Invoice
        output Invoice_db
        err    error
    }{
        "successful post": {
			input:  Invoice{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"},"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1},Item{"Dua","",&itemRate2,itemQuantity2,&itemAmount2}},&itemTax1},
            output:	Invoice_db{Invoice{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"},"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1},Item{"Dua","",&itemRate2,itemQuantity2,&itemAmount2}},&itemTax1},"1903-00001",time.Now().Format("02/01/2006"),&total_1_2_tax1},
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
	srv.PostInvoice(ctx, Invoice{"PT C",nil,"24/03/2018",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax1})
	srv.PostInvoice(ctx, Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax2})
	srv.PostInvoice(ctx, Invoice{"PT A",nil,"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax1})

	tests := map[string]struct {
        input  string
        output Invoice_db
        err    error
    }{
        "successful get": {
            input:  "2003-00001",
            output:	Invoice_db{Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax2},"2003-00001",time.Now().Format("02/01/2006"),&total_1_tax2},
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
	srv.PostInvoice(ctx, Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax2})
	srv.PostInvoice(ctx, Invoice{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"},"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax1})

	tests := map[string]struct {
		input_id	string
        input  		Invoice
        output 		Invoice_db
        err    		error
    }{
        "successful put": {
			input_id:	"1903-00001",
            input:  	Invoice{"PT C",&ClientAddress{Address:"St",PostalCode:""},"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax1},
            output:		Invoice_db{Invoice{"PT C",&ClientAddress{Address:"St",PostalCode:""},"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax1},"1903-00001",time.Now().Format("02/01/2006"),&total_1_tax1},
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
	srv.PostInvoice(ctx, Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax2})
	srv.PostInvoice(ctx, Invoice{"PT A",nil,"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax1})

	tests := map[string]struct {
        input  		string
        output 		bool
        err    		error
    }{
        "successful delete": {
			input:	"1903-00001",
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
	srv.PostInvoice(ctx, Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax2})
	srv.PostInvoice(ctx, Invoice{"PT A",nil,"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax1})

	tests := map[string]struct {
        output map[string]Invoice_db
        err    error
    }{
        "successful get all": {
			output:	map[string]Invoice_db{
				"2003-00001": Invoice_db{Invoice{"PT B",nil,"24/03/2020",[]Item{Item{"Batu","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax2},"2003-00001",time.Now().Format("02/01/2006"),&total_1_tax2},
				"1903-00001": Invoice_db{Invoice{"PT A",nil,"24/03/2019",[]Item{Item{"Paku","",&itemRate1,itemQuantity1,&itemAmount1}},&itemTax1},"1903-00001",time.Now().Format("02/01/2006"),&total_1_tax1},
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