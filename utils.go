package main

import (
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parseData () {

	transaction := Transaction{}
	var (
			newPosting Posting
			commodity string
			amount float64
			price float64
		) 

	reg, err := regexp.Compile(`\d{4}\/(1[0-2]|[1-9])\/(3[0-1]|[1-2][0-9]|[1-9])$`)
	
	if err != nil {
		log.Fatal(err)
		return
	}

	for _,line := range ledgerData {
		transInfo := strings.Split(line, " ")

		if matches := reg.MatchString(transInfo[0]); matches {
			const layout = "2006/1/2"
			date, _ := time.Parse(layout, transInfo[0])
			payee := strings.Join(transInfo[1:]," ")
			

			transactions = append(transactions, transaction)
			transaction = Transaction{}
			transaction.date = date
			transaction.payee = payee

		} else {
			postingInfo := strings.Split(line, "\t")

			account := postingInfo[1]
			quantity := postingInfo[2:]
			if len(quantity) != 0 {
				for _,x := range quantity {
					if x != "" {
						if x[0] == '$' || strings.HasPrefix(x, "-$"){
							commodity = "$"
							price = 1
							amount, _ = strconv.ParseFloat(strings.Replace(x, "$", "", 1), 64)
						} else {
							s := strings.Split(x, " ")
							commodity = s[1]
							price = 0
							amount, _= strconv.ParseFloat(s[0], 64)
						}
						break
					}
				}

			} else {
				amount = amount * (-1)
			}

			newAccount := Account{name: account}

			newPosting = Posting{newAccount, Commodity{name: commodity, price: price}, amount}
			
			newCommodities(commodity, price)
			transaction.postings = append(transaction.postings, newPosting)
			newPosting = Posting{}
		}
	}
	transactions = append(transactions, transaction)
	transactions = transactions[1:]
}

func sortTransactions () {
	sortOption = strings.ToLower(sortOption)
	switch sortOption {
	case "d", "date":
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].date.Before(transactions[j].date)
		})
	case "p", "payee":
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].payee < transactions[j].payee
		})
	case "a", "amount":
		sort.Slice(transactions, func(i, j int) bool {
			return math.Abs(transactions[i].postings[0].amount) < math.Abs(transactions[j].postings[0].amount)
		})
	}
}

func newAccounts(name string, amount float64, comm string) {
	name = strings.TrimSpace(name)
	names := strings.Split(name, ":")

	root.addChildren(names, amount, comm)

}