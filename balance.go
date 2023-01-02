package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)


func orgAccounts () {
	for _, t := range transactions {
		for _, p := range t.postings {
			newAccounts(p.account.name, p.amount, p.commodity.name)
		}
	}
}

func printAccounts(r *Account, argsv []string) {
	if len(r.children) == 1 {
		count := 0
		for kk := range r.children {
			for _,arg := range argsv {
				if arg == r.name {
					for k, v := range r.balance {
					if count > 0 {
						if k == "$" {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%v%.2f", k, v), " ")
						} else {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%.2f %v", v, k), " ")
						}
					} else {
						if k == "$" {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%v%.2f", k, v), r.name+":"+kk)
						} else {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%.2f %v", v, k), r.name+":"+kk)
						}
					}
					count++
				}
				}
			}
			
		}
	} else {
		count := 0

		for _,arg := range argsv {
			if arg == r.name {
				for k, v := range r.balance {
		
					if count > 0 {
						if k == "$" {
						fmt.Printf("%15s %s\n", fmt.Sprintf("%v%.2f", k, v), " ")
						} else {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%.2f %v", v, k), " ")
						}
					} else {
						if k == "$" {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%v%.2f", k, v), r.name)
						} else {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%.2f %v", v, k), r.name)
						}
					}
					count++
				}
			}
		}

		for _,child := range r.children {
			printAccounts(child, argsv)
		}
	}
}

func printAccountsA(r *Account) {
	if len(r.children) == 1 {
		count := 0
		for kk := range r.children {
					for k, v := range r.balance {
					if count > 0 {
						if k == "$" {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%v%.2f", k, v), " ")
						} else {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%.2f %v", v, k), " ")
						}
					} else {
						if k == "$" {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%v%.2f", k, v), r.name+":"+kk)
						} else {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%.2f %v", v, k), r.name+":"+kk)
						}
					}
					count++
			}
			
		}
	} else {
		count := 0
				for k, v := range r.balance {
		
					if count > 0 {
						if k == "$" {
						fmt.Printf("%15s %s\n", fmt.Sprintf("%v%.2f", k, v), " ")
						} else {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%.2f %v", v, k), " ")
						}
					} else {
						if k == "$" {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%v%.2f", k, v), r.name)
						} else {
							fmt.Printf("%15s %s\n", fmt.Sprintf("%.2f %v", v, k), r.name)
						}
					}
					count++
		}

		for _,child := range r.children {
			printAccountsA(child)
		}
	}
}

func balanceCommand (args cli.Args) {
	orgAccounts()
	if args.Len() > 0 {
		printAccounts(root, args.Slice())
	} else {
		printAccountsA(root)
	}

	for _,vC := range root.children {
		for k, v := range vC.balance {
			root.balance[k] += v
		}
	}

	fmt.Println("----------------")
	for k, v := range root.balance {
		if k == "$" {
			fmt.Printf("%15s\n", fmt.Sprintf("%v%.2f", k, v))
		} else {
			fmt.Printf("%15s\n", fmt.Sprintf("%.2f %v", v, k))
		}
	}

}