package report

import (
	"fmt"
	"os"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	gotezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/DefinitelyNotAGoat/go-tezos/delegate"
	"github.com/DefinitelyNotAGoat/payman/cmd/params"
	"github.com/DefinitelyNotAGoat/payman/cmd/print"
	"github.com/spf13/cobra"
)

// NewCommand returns a new report Command
func NewCommand() *cobra.Command {
	var confFile string
	var cycle int
	var report = &cobra.Command{
		Use:   "report",
		Short: "report simulates a payout",
		Run: func(cmd *cobra.Command, args []string) {

			conf, err := params.NewReportParameters(confFile)
			if err != nil {
				fmt.Println("Missing required parameters.")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			gt, err := goTezos.NewGoTezos(conf.URL)
			if err != nil {
				fmt.Printf("Could not connect to Tezos node: %s\n", conf.URL)
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			if cycle == 0 {
				cycle = conf.Cycle
			}

			report, err := gt.Delegate.GetReport(
				conf.Delegate,
				cycle,
				float64(conf.Fee),
				false,
			)
			if err != nil {
				fmt.Printf("Could not generate a report\n")
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			print.ReportTable(*report)
		},
	}

	report.PersistentFlags().StringVarP(&confFile, "conf", "c", "", "overide the default config file (~/.payman/report.json)")
	report.PersistentFlags().IntVar(&cycle, "cycle", 0, "overide the default cycle from the config")

	return report
}

func report(conf params.PayoutParameters, gt *gotezos.GoTezos) (*delegate.DelegateReport, error) {
	report, err := gt.Delegate.GetReport(
		conf.Delegate,
		conf.Cycle,
		float64(conf.Fee),
		false,
	)
	if err != nil {
		return report, err
	}

	return report, nil
}
