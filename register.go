package main

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

func printRegisterTable (table []registerRow) {
	for _,row := range table {
		if len(row.payee) > 30 {
			row.payee = row.payee[:30] + ".."
		}
		if len(row.account) > 20 {
			row.account = row.account[:20] + ".."
		}
		fmt.Printf("%10s %-35s %-25s %15s %15s\n", row.date, row.payee, row.account, row.amount, row.rBalance)
	}
}

func registerCommand (args cli.Args) {

	registerTable := []registerRow{}
	runningBalances := make(map[string]float64)
	var rowAux registerRow

	for _,transaction := range transactions {

		//fmt.Printf("\n%10s %-30s\t", time.Time.Format(transaction.date, "06-Jan-02"), transaction.payee)
		rowAux.date = time.Time.Format(transaction.date, "06-Jan-02")
		rowAux.payee = transaction.payee

		for _,posting := range transaction.postings {
			runningBalances[posting.commodity.name] += posting.amount
			rowAux.account = posting.account.name

			if posting.commodity.name == "$" {
				rowAux.amount = fmt.Sprintf("%v%.2f", posting.commodity.name, posting.amount)
			} else {
				rowAux.amount = fmt.Sprintf("%.1f %v", posting.amount, posting.commodity.name)
			}

			for commodity, balance := range runningBalances {
				if commodity == "$" {
					rowAux.rBalance = fmt.Sprintf("%v%.2f", commodity, balance)
				} else {
					rowAux.rBalance = fmt.Sprintf("%.1f %v", balance, commodity)
				}
				registerTable = append(registerTable, rowAux)
				rowAux = registerRow{}
			}
		}
	}
	printRegisterTable(registerTable)
}