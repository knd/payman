package payout

import (
	"fmt"
	"log"
	"os"
	"strconv"

	gotezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/DefinitelyNotAGoat/go-tezos/account"
	"github.com/DefinitelyNotAGoat/go-tezos/delegate"
	"github.com/DefinitelyNotAGoat/payman/twitter"

	"github.com/DefinitelyNotAGoat/payman/cmd/batch"
	"github.com/DefinitelyNotAGoat/payman/cmd/params"
	"github.com/DefinitelyNotAGoat/payman/cmd/print"

	"github.com/spf13/cobra"
)

// NewCommand returns a new payout command
func NewCommand() *cobra.Command {
	var confFile string
	var csv bool
	var payout = &cobra.Command{
		Use:   "payout",
		Short: "payout pays out rewards to delegations.",
		Long:  "payout computes rewards for a delegates delegations and batch pays them based off payout's configurable parameters.",
		Run: func(cmd *cobra.Command, args []string) {

			conf, err := params.NewPayoutParameters(confFile)
			if err != nil {
				fmt.Println("Missing required parameters.")
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			gt, err := gotezos.NewGoTezos(conf.URL)
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

			var twitterBot *twitter.Bot
			if conf.Twitter != nil {
				twitterBot, err = twitter.NewBot(
					"",
					conf.Twitter.ConsumerKey,
					conf.Twitter.ConsumerKeySecret,
					conf.Twitter.AccessToken,
					conf.Twitter.AccessTokenSecret,
				)
				if err != nil {
					fmt.Println("Could not start twitter bot")
					fmt.Printf("Error: %s\n", err.Error())
					os.Exit(1)
				}
			}

			result, hashes, err := payout(conf, wallet, gt)
			if err != nil {
				log.Fatal(err)
			}

			for _, hash := range hashes {
				fmt.Printf("Successful operation: %s\n", hash)
				if conf.Twitter != nil {
					err := twitterBot.Post(hash, conf.Cycle)
					if err != nil {
						fmt.Println("Could not post to twitter")
						fmt.Printf("Error: %s\n", err.Error())
					}
				}
			}

			print.ReportTable(*result)
			if csv {
				print.WriteCSV(*result)
			}

		},
	}
	payout.PersistentFlags().StringVarP(&confFile, "conf", "c", "", "overide the default config file (~/.payman/payouts.json)")
	payout.PersistentFlags().BoolVar(&csv, "csv", false, "writes a csv version of the result of payout")
	return payout
}

func payout(conf params.PayoutParameters, wallet account.Wallet, gt *gotezos.GoTezos) (*delegate.DelegateReport, []string, error) {
	report, err := gt.Delegate.GetReport(
		conf.Delegate,
		conf.Cycle,
		float64(conf.Fee),
		false,
	)
	if err != nil {
		return report, []string{}, err
	}

	payments := report.GetPayments(conf.PaymentMinimum)
	report.Delegations = cleanseReport(conf.PaymentMinimum, report.Delegations)

	hashes, err := batch.Batch(conf.NetworkFee, conf.NetworkGasLimit, wallet, payments, gt)
	if err != nil {
		return report, []string{}, err
	}

	return report, hashes, nil
}

func cleanseReport(payMin int, delegations []delegate.DelegationReport) []delegate.DelegationReport {
	for _, delegation := range delegations {
		intNet, _ := strconv.Atoi(delegation.NetRewards)
		if intNet >= payMin {
			delegations = append(delegations, delegation)
		}
	}

	return delegations
}
