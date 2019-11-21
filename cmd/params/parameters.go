package params

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/go-playground/validator"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

var validate = validator.New()

const (
	// BATCHFILE default batch configuration
	BATCHFILE = ".payman/batch.json"
	// PAYOUTFILE default payout configuration
	PAYOUTFILE = ".payman/payout.json"
	// REPORTFILE default report configuration
	REPORTFILE = ".payman/report.json"
)

// PayoutParameters is a configuration to payman payout
type PayoutParameters struct {
	Delegate        string             `validate:"required" json:"delegate"`
	Cycle           int                `validate:"required" json:"cycle"`
	Fee             float32            `validate:"required" json:"fee"`
	Wallet          WalletParameters   `validate:"required" json:"wallet"`
	URL             string             `default:"http://127.0.0.1:8732" json:"url"`
	NetworkFee      int                `default:"2941" json:"network_fee"`
	NetworkGasLimit int                `default:"26283" json:"gas_limit"`
	PaymentMinimum  int                `json:"payment_min"`
	Twitter         *TwitterParameters `json:"twitter"`
}

// WalletParameters is the wallet to payout with
type WalletParameters struct {
	Secret   string `validate:"required" json:"secret"`
	Password string `validate:"required" json:"password"`
}

// TwitterParameters is the api information needed to tweet a payout
type TwitterParameters struct {
	ConsumerKey       string `validate:"required" json:"key"`
	ConsumerKeySecret string `validate:"required" json:"secret"`
	AccessToken       string `validate:"required" json:"access_token"`
	AccessTokenSecret string `validate:"required" json:"access_secret"`
}

// ReportParameters is a configuration to payman report
type ReportParameters struct {
	Delegate       string  `validate:"required" json:"delegate"`
	Cycle          int     `validate:"required" json:"cycle"`
	Fee            float32 `validate:"required" json:"fee"`
	URL            string  `default:"http://127.0.0.1:8732" json:"url"`
	PaymentMinimum int     `json:"payment_min"`
}

//BatchParamerters is a configuration for batch payments
type BatchParamerters struct {
	Wallet          WalletParameters `validate:"required" json:"wallet"`
	URL             string           `default:"http://127.0.0.1:8732" json:"url"`
	NetworkFee      int              `default:"2941" json:"network_fee"`
	NetworkGasLimit int              `default:"26283" json:"gas_limit"`
}

// NewPayoutParameters will read in payman payout parametrs from the passed file.
// If the file doesn't exist, enviroment variables will be tried.
func NewPayoutParameters(paramsFile string) (PayoutParameters, error) {
	var payoutParams PayoutParameters
	if err := envconfig.Process("payman", &payoutParams); err != nil {
		return payoutParams, errors.New("failed to get payout parameters")
	}
	if validate.Struct(payoutParams) != nil {
		if paramsFile == "" {
			usr, err := user.Current()
			if err != nil {
				return payoutParams, err
			}
			paramsFile = fmt.Sprintf("%s/%s", usr.HomeDir, PAYOUTFILE)
		}

		byts, err := readFile(paramsFile)
		if err != nil {
			return payoutParams, err
		}
		err = json.Unmarshal(byts, &payoutParams)
		if err != nil {
			return payoutParams, err
		}
	}

	return payoutParams, validate.Struct(payoutParams)
}

// NewReportParameters will read in payman report parametrs from the passed file.
// If the file doesn't exist, enviroment variables will be tried.
func NewReportParameters(paramsFile string) (ReportParameters, error) {
	var reportParams ReportParameters
	if err := envconfig.Process("payman", &reportParams); err != nil {
		return reportParams, errors.New("failed to get report parameters")
	}
	if validate.Struct(reportParams) != nil {
		if paramsFile == "" {
			usr, err := user.Current()
			if err != nil {
				return reportParams, err
			}
			paramsFile = fmt.Sprintf("%s/%s", usr.HomeDir, REPORTFILE)
		}

		byts, err := readFile(paramsFile)
		if err != nil {
			return reportParams, err
		}
		err = json.Unmarshal(byts, &reportParams)
		if err != nil {
			return reportParams, err
		}
	}

	return reportParams, validate.Struct(reportParams)
}

// NewBatchParameters will read in payman batch parametrs from the passed file.
// If the file doesn't exist, enviroment variables will be tried.
func NewBatchParameters(paramsFile string) (BatchParamerters, error) {
	var batchParams BatchParamerters

	if err := envconfig.Process("payman", &batchParams); err != nil {
		return batchParams, errors.New("failed to get batch parameters")
	}
	if validate.Struct(batchParams) != nil {
		if paramsFile == "" {
			usr, err := user.Current()
			if err != nil {
				return batchParams, err
			}
			paramsFile = fmt.Sprintf("%s/%s", usr.HomeDir, BATCHFILE)
		}

		byts, err := readFile(paramsFile)
		if err != nil {
			return batchParams, err
		}
		err = json.Unmarshal(byts, &batchParams)
		if err != nil {
			return batchParams, err
		}
	}

	return batchParams, validate.Struct(batchParams)
}

func readFile(file string) ([]byte, error) {
	pFile, err := os.Open(file)
	if err != nil {
		return []byte{}, err
	}
	defer pFile.Close()
	byteValue, err := ioutil.ReadAll(pFile)
	if err != nil {
		return byteValue, err
	}

	return byteValue, err
}
