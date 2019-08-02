package options

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/DefinitelyNotAGoat/go-tezos/delegate"
)

// Options is a struct to represent configuration options for payman
type Options struct {
	Delegate         string
	Secret           string
	Password         string
	Service          bool
	Cycle            int
	URL              string
	Fee              float32
	File             string
	NetworkFee       int
	NetworkGasLimit  int
	PaymentMinimum   int
	Dry              bool
	RedditAgent      string
	RedditTitle      string
	TwitterPath      string
	TwitterTitle     string
	Twitter          bool
	PaymentsOverride PaymentsOverride
}

//PaymentsOverride is a configuration option to override the payments calculation with your own
type PaymentsOverride struct {
	File     string
	Payments []delegate.Payment
}

// ReadPaymentsOverride generates the Payment Struct for the payer
func (p *PaymentsOverride) ReadPaymentsOverride() ([]delegate.Payment, error) {
	jsonFile, err := os.Open(p.File)
	if err != nil {
		return []delegate.Payment{}, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []delegate.Payment{}, err
	}

	var payments []delegate.Payment
	err = json.Unmarshal(byteValue, &payments)
	if err != nil {
		return payments, err
	}

	return payments, nil
}
