package main

import (
	"fmt"
	"math"
	"time"
)

func printCommand () {
	for _, transaction := range transactions {
		fmt.Printf("\n%v %v \n", time.Time.Format(transaction.date, "2006/01/02"), transaction.payee)
		if len(transaction.postings) == 2 {
			if transaction.postings[0].commodity.name ==  transaction.postings[1].commodity.name{
				if math.Abs(transaction.postings[0].amount) == math.Abs(transaction.postings[1].amount) {
					if transaction.postings[0].commodity.name == "$" {
						fmt.Printf("    %-30s %15s\n", transaction.postings[0].account.name, fmt.Sprintf("%v%.2f", transaction.postings[0].commodity.name, transaction.postings[0].amount))
						fmt.Printf("    %-30s\n", transaction.postings[1].account.name)
					} else {
						fmt.Printf("    %-30s %15s\n", transaction.postings[0].account.name, fmt.Sprintf("%.1f %v", transaction.postings[0].amount, transaction.postings[0].commodity.name))
						fmt.Printf("    %-30s\n", transaction.postings[1].account.name)
					}
				} else {
					if transaction.postings[0].commodity.name == "$" {
						fmt.Printf("    %-30s %15s\n", transaction.postings[0].account.name, fmt.Sprintf("%v%.2f", transaction.postings[0].commodity.name, transaction.postings[0].amount))
					} else {
						fmt.Printf("    %-30s %15s\n", transaction.postings[0].account.name, fmt.Sprintf("%.1f %v", transaction.postings[0].amount, transaction.postings[0].commodity.name))
					}
					if transaction.postings[1].commodity.name == "$" {
						fmt.Printf("    %-30s %15s\n", transaction.postings[1].account.name, fmt.Sprintf("%v%.2f", transaction.postings[1].commodity.name, transaction.postings[1].amount))
					} else {
						fmt.Printf("    %-30s %15s\n", transaction.postings[1].account.name, fmt.Sprintf("%.1f %v", transaction.postings[1].amount, transaction.postings[1].commodity.name))
					}
				}
			} else {
				if transaction.postings[0].commodity.name == "$" {
					transaction.postings[1].commodity.price = transaction.postings[0].amount / transaction.postings[1].amount
					fmt.Printf("    %-30s %15s\n", transaction.postings[0].account.name, fmt.Sprintf("%v%.2f", transaction.postings[0].commodity.name, transaction.postings[0].amount))
				} else {
					fmt.Printf("    %-30s %15s\n", transaction.postings[0].account.name, fmt.Sprintf("%.1f %v", transaction.postings[0].amount, transaction.postings[0].commodity.name))
				}
				if transaction.postings[1].commodity.name == "$" {
					transaction.postings[0].commodity.price = transaction.postings[1].amount / transaction.postings[0].amount
					fmt.Printf("    %-30s %15s\n", transaction.postings[1].account.name, fmt.Sprintf("%v%.2f", transaction.postings[1].commodity.name, transaction.postings[1].amount))
				} else {
					fmt.Printf("    %-30s %15s\n", transaction.postings[1].account.name, fmt.Sprintf("%.1f %v", transaction.postings[1].amount, transaction.postings[1].commodity.name))
				}
			}
		} else {
			for _, posting := range transaction.postings {
				if posting.commodity.name == "$" {
					fmt.Printf("    %-30s %15s\n", posting.account.name, fmt.Sprintf("%v%.2f", posting.commodity.name, posting.amount))
				} else {
					fmt.Printf("    %-30s %15s\n", posting.account.name, fmt.Sprintf("%.1f %v", posting.amount, posting.commodity.name))
				}
			}
		}
	}
}