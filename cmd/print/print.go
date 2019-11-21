package print

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"encoding/csv"

	"github.com/DefinitelyNotAGoat/go-tezos/account"
	"github.com/DefinitelyNotAGoat/go-tezos/delegate"
	"github.com/olekukonko/tablewriter"
)

// ReportTable takes in payments and prints them to a table for general logging
func ReportTable(report delegate.DelegateReport) {
	total := []string{}
	data := formatData(report)
	if len(data) > 0 {
		total = data[len(data)-1]
		data = data[:len(data)-1]
	}

	table := tablewriter.NewWriter(&writer{})
	table.SetHeader([]string{"Address", "Share", "Gross", "Fee", "Net"})
	table.SetFooter(total)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// formatData parses payments into a double array of data for table or csv printing
func formatData(report delegate.DelegateReport) [][]string {
	var data [][]string
	var totalNet float64
	var totalGross float64
	var totalFee float64
	for _, payment := range report.Delegations {
		share := payment.Share * 100
		strShare := fmt.Sprintf("%.6f", share)
		fee, _ := strconv.Atoi(payment.Fee)
		floatFee := float64(fee) / float64(account.MUTEZ)
		gross, _ := strconv.Atoi(payment.GrossRewards)
		floatGross := float64(gross) / float64(account.MUTEZ)
		net, _ := strconv.Atoi(payment.NetRewards)
		floatNet := float64(net) / float64(account.MUTEZ)

		totalNet = totalNet + floatNet
		totalGross = totalGross + floatGross
		totalFee = totalFee + floatFee
		data = append(data, []string{payment.DelegationPhk, strShare, fmt.Sprintf("%.6f", floatGross), fmt.Sprintf("%.6f", floatFee), fmt.Sprintf("%.6f", floatNet)})
	}
	data = append(data, []string{"", "Total", fmt.Sprintf("%.6f", totalGross), fmt.Sprintf("%.6f", totalFee), fmt.Sprintf("%.6f", totalNet)})
	return data
}

// WriteCSV writes payments to a csv file for reporting
func WriteCSV(payments delegate.DelegateReport) error {
	r, err := getCSVWriter()
	if err != nil {
		return err
	}

	data := formatData(payments)
	for _, value := range data {
		r.Write(value)
	}
	r.Flush()
	return nil
}

// getCSVWriter opens a file with the current date for its name and
// returns a csv.Writer for that file
func getCSVWriter() (*csv.Writer, error) {
	fileName := buildFileName()
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	report := csv.NewWriter(f)
	return report, nil
}

// buildFileName returns a string that represents a filename based
// of the current date
func buildFileName() string {
	return time.Now().Format(time.RFC3339) + ".csv"
}

type writer struct{}

func (w *writer) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}
