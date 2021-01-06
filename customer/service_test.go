package customer

import (
  "context"
	"testing"
	"time"
	"strconv"
	"github.com/stretchr/testify/assert"
  "github.com/steven-steven/GoInvoice/config"
)

var phone = "12345"
var id1 = "customer_1"
var id2 = "customer_2"

func mock_defer_idGenerator() func() {
	// mock idGenerator() to return id1.
	origIdGenerator := idGenerator
	num := 0
	idGenerator = func() string {
		num += 1
		return "customer_" + strconv.Itoa(num)
	}
	return func() { idGenerator = origIdGenerator }
}

func TestPostCustomer(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup()	//new test DB
	
	tests := map[string]struct {
			input  Customer
			output Customer_db
			err    error
    }{
			"successful post": {
			input:  Customer{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"}, phone},
			output:	Customer_db{Customer{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"}, phone},id1,time.Now().Format("02/01/2006")},
			err:    nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.PostCustomer(ctx, test.input)
			assert.IsType(t, test.err, err)
        	assert.EqualValues(t, test.output, output)
		})
	}
}

func TestGetCustomer(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostCustomer(ctx, Customer{"PT C",nil, phone})
	srv.PostCustomer(ctx, Customer{"PT B",nil, phone})
	srv.PostCustomer(ctx, Customer{"PT A",nil, phone})

	tests := map[string]struct {
			input  string
			output Customer_db
			err    error
    }{
			"successful get": {
				input:  id2,
				output:	Customer_db{Customer{"PT B",nil, phone},id2,time.Now().Format("02/01/2006")},
				err:    nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.GetCustomer(ctx, test.input)
			assert.IsType(t, test.err, err)
        	assert.EqualValues(t, test.output, output)
		})
	}
}

func TestPutCustomer(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostCustomer(ctx, Customer{"PT B",nil,"123"})
	srv.PostCustomer(ctx, Customer{"PT A",&ClientAddress{"690 King St","Cilegon","Banten","Indonesia","154321"}, phone})

	tests := map[string]struct {
			input_id	string
			input  		Customer
			output 		Customer_db
			err    		error
    }{
			"successful put": {
				input_id:	id2,
				input:  	Customer{"PT C",&ClientAddress{Address:"St",PostalCode:""}, phone},
				output:		Customer_db{Customer{"PT C",&ClientAddress{Address:"St",PostalCode:""}, phone},id2,time.Now().Format("02/01/2006")},
				err:    	nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.PutCustomer(ctx, test.input_id, test.input)
			assert.IsType(t, test.err, err)
			assert.EqualValues(t, test.output, output)
			//check get returns updated
			after, errAfter := srv.GetCustomer(ctx, test.input_id)
			assert.IsType(t, nil, errAfter)
			assert.EqualValues(t, test.output, after)
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostCustomer(ctx, Customer{"PT B",nil, phone})
	srv.PostCustomer(ctx, Customer{"PT A",nil, phone})

	tests := map[string]struct {
			input  		string
			output 		bool
			err    		error
    }{
      "successful delete": {
				input:	id2,
				output: true,
				err:	nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.DeleteCustomer(ctx, test.input)
			assert.IsType(t, test.err, err)
			assert.EqualValues(t, test.output, output)
			//check get returns updated
			after, errAfter := srv.GetCustomer(ctx, test.input)
			assert.IsType(t, ApiError, errAfter)
			assert.EqualValues(t, Customer_db{}, after)
		})
	}
}

func TestGetAllCustomer(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostCustomer(ctx, Customer{"PT B",nil, phone})
	srv.PostCustomer(ctx, Customer{"PT A",nil, phone})

	tests := map[string]struct {
			output map[string]Customer_db
			err    error
    }{
			"successful get all": {
				output:	map[string]Customer_db{
				id1: Customer_db{Customer{"PT B",nil, phone},id1,time.Now().Format("02/01/2006")},
				id2: Customer_db{Customer{"PT A",nil, phone},id2,time.Now().Format("02/01/2006")},
			},
			err:    nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.GetAllCustomer(ctx)
			assert.IsType(t, test.err, err)
        	assert.EqualValues(t, test.output, output)
		})
	}
}

func setup() (srv Service, ctx context.Context) {
	dbClient := config.GetTestDB(ctx)
  return NewService(dbClient), context.Background()
}