package batch

import (
	"fmt"
	"os"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	gotezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/DefinitelyNotAGoat/go-tezos/account"
	"github.com/DefinitelyNotAGoat/go-tezos/delegate"
	"github.com/DefinitelyNotAGoat/payman/cmd/params"
	"github.com/spf13/cobra"
)

// NewCommand returns a new batch Command
func NewCommand() *cobra.Command {
	var confFile string
	var batch = &cobra.Command{
		Use:   "batch",
		Short: "batch performs batch/bulk transfers",
		Run: func(cmd *cobra.Command, args []string) {

			conf, err := params.NewBatchParameters(confFile)
			if err != nil {
				fmt.Println("Missing required parameters.")
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			gt, err := goTezos.NewGoTezos(conf.URL)
			if err != nil {
				fmt.Printf("Could not connect to Tezos node: %s\n", conf.URL)
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			wallet, err := gt.Account.ImportEncryptedWallet(conf.Wallet.Password, conf.Wallet.Secret)
			if err != nil {
				fmt.Println("Could not import wallet.")
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			hashes, err := Batch(conf.NetworkFee, conf.NetworkGasLimit, wallet, nil, gt)
			if err != nil {
				fmt.Println("Could make batch payment.")
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			for _, hash := range hashes {
				fmt.Printf("Batch Payment Successful: https://mvp.tezblock.io/transaction/%s", hash)
			}

		},
	}

	batch.PersistentFlags().StringVarP(&confFile, "conf", "c", "", "overide the default config file (~/.payman/report.json)")

	return batch
}

// Batch makes a batch of transfers and returns the operations hashes
func Batch(networkFee, gasLimit int, wallet account.Wallet, payments []delegate.Payment, gt *gotezos.GoTezos) ([]string, error) {
	ops, err := gt.Operation.CreateBatchPayment(
		payments,
		wallet,
		networkFee,
		gasLimit,
		100,
	)
	if err != nil {
		return []string{}, err
	}

	hashes, err := injectOperations(gt, ops...)
	if err != nil {
		return hashes, err
	}

	return hashes, nil
}

// injectOperations injects operation strings and returns their hashes
func injectOperations(gt *gotezos.GoTezos, ops ...string) ([]string, error) {
	responses := [][]byte{}
	for _, op := range ops {
		resp, err := gt.Operation.InjectOperation(op)
		if err != nil {
			return []string{}, err
		}
		responses = append(responses, resp)
	}

	var strHashes []string
	for _, hash := range responses {
		strHashes = append(strHashes, string(hash))
	}

	return strHashes, nil
}
