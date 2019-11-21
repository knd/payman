package cmd

import (
	"github.com/DefinitelyNotAGoat/payman/cmd/batch"
	"github.com/DefinitelyNotAGoat/payman/cmd/payout"
	"github.com/DefinitelyNotAGoat/payman/cmd/report"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "payman",
		Short: "A bulk payout tool for bakers in the Tezos Ecosystem",
	}

	rootCommand.AddCommand(
		batch.NewCommand(),
		payout.NewCommand(),
		report.NewCommand(),
	)

	return rootCommand
}

func Execute() error {
	return newRootCommand().Execute()
}
