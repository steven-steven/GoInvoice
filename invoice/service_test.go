package invoice

import (
	"context"
	"strconv"
	"testing"
	"time"

	// "reflect"
	"math"

	"github.com/steven-steven/GoInvoice/config"
	"github.com/stretchr/testify/assert"
)

var itemRate1 = uint64(10000)
var itemQuantity1 = uint64(3000)
var itemAmount1 = itemRate1 * uint64(itemQuantity1) / 1000
var itemRate2 = uint64(5000)
var itemQuantity2 = uint64(543500)
var itemAmount2 = uint64(float64(itemRate2) * float64(itemQuantity2) / 1000)
var unit1 = "m^2"
var unit2 = "unit"
var itemTax1 = uint64(5000)
var itemTax2 = uint64(6000)

var subtotal_1_2 = itemAmount1 + itemAmount2
var total_1_2_tax1 = subtotal_1_2 + uint64(math.Round(float64(itemTax1)/100*float64(subtotal_1_2)))
var total_1_2_tax2 = subtotal_1_2 + uint64(math.Round(float64(itemTax2)/100*float64(subtotal_1_2)))

var subtotal_1 = itemAmount1
var subtotal_2 = itemAmount2
var total_1_tax1 = itemAmount1 + uint64(math.Round(float64(itemTax1)/100*float64(itemAmount1)))
var total_2_tax2 = itemAmount2 + uint64(math.Round(float64(itemTax2)/100*float64(itemAmount2)))
var total_2_tax1 = itemAmount2 + uint64(math.Round(float64(itemTax1)/100*float64(itemAmount2)))
var total_1_tax2 = itemAmount1 + uint64(math.Round(float64(itemTax2)/100*float64(itemAmount1)))

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

func TestPostInvoice(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup() //new test DB

	tests := map[string]struct {
		input  Invoice
		output Invoice_db
		err    error
	}{
		"successful post": {
			input:  Invoice{"invNo1", "PT A", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{Item{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}, Item{"Dua", "", &itemRate2, true, unit1, &itemQuantity2, &itemAmount2}}, &itemTax1, false},
			output: Invoice_db{Invoice{"invNo1", "PT A", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{Item{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}, Item{"Dua", "", &itemRate2, true, unit1, &itemQuantity2, &itemAmount2}}, &itemTax1, false}, id1, time.Now().Format("02/01/2006"), &total_1_2_tax1, &subtotal_1_2},
			err:    nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T) {
			output, err := srv.PostInvoice(ctx, test.input)
			assert.IsType(t, test.err, err)
			assert.EqualValues(t, test.output, output)
		})
	}
}

func TestGetInvoice(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup() //new test DB

	//initial data
	srv.PostInvoice(ctx, Invoice{"invNo1", "PT C", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2018", []Item{Item{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax1, false})
	srv.PostInvoice(ctx, Invoice{"invNo2", "PT B", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2020", []Item{Item{"Batu", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax2, false})
	srv.PostInvoice(ctx, Invoice{"invNo3", "PT A", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{Item{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax1, false})

	tests := map[string]struct {
		input  string
		output Invoice_db
		err    error
	}{
		"successful get": {
			input:  id2,
			output: Invoice_db{Invoice{"invNo2", "PT B", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2020", []Item{Item{"Batu", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax2, false}, id2, time.Now().Format("02/01/2006"), &total_1_tax2, &subtotal_1},
			err:    nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T) {
			output, err := srv.GetInvoice(ctx, test.input)
			assert.IsType(t, test.err, err)
			assert.EqualValues(t, test.output, output)
		})
	}
}

func TestPutInvoice(t *testing.T) {
	defer mock_defer_idGenerator()()
	srv, ctx := setup() //new test DB

	//initial data
	srv.PostInvoice(ctx, Invoice{"invNo1", "PT B", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2020", []Item{Item{"Batu", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax2, false})
	srv.PostInvoice(ctx, Invoice{"invNo2", "PT A", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{Item{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax1, false})

	tests := map[string]struct {
		input_id string
		input    Invoice
		output   Invoice_db
		err      error
	}{
		"successful put": {
			input_id: id2,
			input:    Invoice{"invNo2", "PT C", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{Item{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax1, true},
			output:   Invoice_db{Invoice{"invNo2", "PT C", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{Item{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax1, true}, id2, time.Now().Format("02/01/2006"), &total_1_tax1, &subtotal_1},
			err:      nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T) {
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
	defer mock_defer_idGenerator()()
	srv, ctx := setup() //new test DB

	//initial data
	srv.PostInvoice(ctx, Invoice{"invNo1", "PT B", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2020", []Item{Item{"Batu", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax2, false})
	srv.PostInvoice(ctx, Invoice{"invNo2", "PT A", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{Item{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax1, false})

	tests := map[string]struct {
		input  string
		output bool
		err    error
	}{
		"successful delete": {
			input:  id2,
			output: true,
			err:    nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T) {
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
	defer mock_defer_idGenerator()()
	srv, ctx := setup() //new test DB

	//initial data
	srv.PostInvoice(ctx, Invoice{"invNo1", "PT B", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2020", []Item{{"Batu", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax2, false})
	srv.PostInvoice(ctx, Invoice{"invNo2", "PT A", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax1, false})

	tests := map[string]struct {
		output map[string]Invoice_db
		err    error
	}{
		"successful get all": {
			output: map[string]Invoice_db{
				id1: {Invoice{"invNo1", "PT B", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2020", []Item{{"Batu", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax2, false}, id1, time.Now().Format("02/01/2006"), &total_1_tax2, &subtotal_1},
				id2: {Invoice{"invNo2", "PT A", "catatanInvoice", "catatanKwi", "keteranganKwitansi", "24/03/2019", []Item{{"Paku", "", &itemRate1, false, unit2, &itemQuantity1, &itemAmount1}}, &itemTax1, false}, id2, time.Now().Format("02/01/2006"), &total_1_tax1, &subtotal_1},
			},
			err: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		t.Run(testName, func(t *testing.T) {
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
