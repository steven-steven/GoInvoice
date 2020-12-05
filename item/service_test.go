package item

import (
    "context"
    "testing"
	"time"
	"strconv"
	// "reflect"
	"github.com/stretchr/testify/assert"
	"github.com/steven-steven/GoInvoice/config"
)

var itemRate1 = uint64(10000)
var id1 = "item_1"
var id2 = "item_2"

func mock_defer_idGenerator() func() {
	// mock idGenerator() to return id1.
	origIdGenerator := idGenerator
	num := 0
	idGenerator = func() string {
		num += 1
		return "item_" + strconv.Itoa(num)
	}
	return func() { idGenerator = origIdGenerator }
}

func TestPostItem(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup()	//new test DB
	
	tests := map[string]struct {
        input  Item
        output Item_db
        err    error
    }{
        "successful post": {
			input:  Item{"PT A","defaultDescription",&itemRate1},
            output:	Item_db{Item{"PT A","defaultDescription",&itemRate1},id1,time.Now().Format("02/01/2006")},
           	err:    nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.PostItem(ctx, test.input)
			assert.IsType(t, test.err, err)
        	assert.EqualValues(t, test.output, output)
		})
	}
}

func TestDeleteItem(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostItem(ctx, Item{"PT B","defaultDescription",&itemRate1})
	srv.PostItem(ctx, Item{"PT A","defaultDescription",&itemRate1})

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
			output, err := srv.DeleteItem(ctx, test.input)
			assert.IsType(t, test.err, err)
			assert.EqualValues(t, test.output, output)
			//check get returns updated
			after, errAfter := srv.GetItem(ctx, test.input)
			assert.IsType(t, ApiError, errAfter)
			assert.EqualValues(t, Item_db{}, after)
		})
	}
}

func TestGetAllItem(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup()	//new test DB
	
	//initial data
	srv.PostItem(ctx, Item{"PT B","defaultDescription",&itemRate1})
	srv.PostItem(ctx, Item{"PT A","defaultDescription",&itemRate1})

	tests := map[string]struct {
        output map[string]Item_db
        err    error
    }{
        "successful get all": {
			output:	map[string]Item_db{
				id1: Item_db{Item{"PT B","defaultDescription",&itemRate1},id1,time.Now().Format("02/01/2006")},
				id2: Item_db{Item{"PT A","defaultDescription",&itemRate1},id2,time.Now().Format("02/01/2006")},
			},
           	err:    nil,
		},
	}
	
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T){
			output, err := srv.GetAllItem(ctx)
			assert.IsType(t, test.err, err)
        	assert.EqualValues(t, test.output, output)
		})
	}
}

func setup() (srv Service, ctx context.Context) {
	dbClient := config.GetTestDB(ctx)
    return NewService(dbClient), context.Background()
}