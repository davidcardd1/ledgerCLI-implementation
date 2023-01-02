package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

var (
	app = cli.NewApp()
	file string
	sortOption string
	priceDBFile string
	comments = []rune{'!', ';', '#', '%', '|', '*'}
	ledgerData = []string{}
	transactions = []Transaction{}
	accounts = []Account{}
	commodities = []Commodity{}
)

type Transaction struct {
	date 		time.Time
	payee 		string
	postings	[]Posting
}

type Posting struct {
	account 		Account
	commodity 		Commodity
	amount 			float64
}

type Account struct {
	name 		string
	hasposting	bool
}

type Commodity struct {
	name 	string
	price	float64 
}

type registerRow struct {
		date		string
		payee		string
		account		string
		amount		string
		rBalance	string
	}

func fileReader(file string) {
	f, err := os.Open("./ledger-sample-files/"+file)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}


	for _, line := range lines {
		isComment := false
		if strings.HasPrefix(line, "!include") {
			fileReader(strings.Split(line, " ")[1])
		}
		for _, c := range comments {
			if c == rune(line[0]) {
				isComment = true
				break
			}
		}
		if !isComment {
			// fmt.Println(line)
			ledgerData = append(ledgerData, line)
		}
	}
}

func appendTransaction(transactions *[]Transaction, transaction *Transaction) {
	var sum float64
	for _,posting := range transaction.postings {
		sum += posting.amount * posting.commodity.price
	}

	if sum > 0.05 || sum <= -.05 {
		log.Fatal("Transaction doesn't balance");
	} else {
		*transactions = append(*transactions, *transaction)
	}
}

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

				newAccount := Account{account, true}
				newCommodity := Commodity{name: commodity, price: price}

				accounts = append(accounts, newAccount)
				commodities = append(commodities, newCommodity)

				newPosting = Posting{newAccount, newCommodity, amount}

			} else {
				newAccount := Account{account, true}
				accounts = append(accounts, newAccount)

				newPosting = Posting{newAccount, Commodity{name: commodity, price: price}, -amount}

			}
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

func printRegisterTable (table []registerRow) {
	for _,row := range table {
		if len(row.payee) > 30 {
			row.payee = row.payee[:29] + ".."
		}
		if len(row.account) > 20 {
			row.account = row.account[:19] + ".."
		}
		fmt.Printf("%10s %-35s %-25s %15s %15s\n", row.date, row.payee, row.account, row.amount, row.rBalance)
	}
}

func registerCommand () {

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

func balanceCommand () {

}

func flags() {
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "sort",
			Aliases: []string{"s", "S"},
			Value:   "date",
			Usage:   "Sort report using `VEXPR` which is either of ['date'/'d', 'payee'/'p', 'amount'/'a']",
			Destination: &sortOption,
			Action: func(ctx *cli.Context, s string) error {
				options := []string{"date", "d", "payee", "p", "amount", "a"}
				if !slices.Contains(options, s) {
					log.Fatalf("Flag sort value '%v' is not 'date', 'payee' or 'amount'", s)
				}				
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Read ledger file using `FILE`",
			Value: "index.ledger",
			Destination: &file,
		},
		&cli.StringFlag{
			Name:    "price-db",
			Usage:   "Use `FILE` for retrieving stored commodity prices.",
			Destination: &priceDBFile,
		},
	}
}

func commands() {
	app.Commands = []*cli.Command{
		{
			Name: "balance",
			Aliases: []string{"bal"},
			Usage: "Current balance of all accounts, aggregating totals for parent accounts and different commodities",
			Action: func (c *cli.Context) error {
				fmt.Println("selected balance")
				fileReader(file)
				parseData()
				sortTransactions()
				balanceCommand()
				return nil
			},
		},
		{
			Name: "register",
			Usage: "Lists all postings that match the report-query with running total",
			Aliases: []string{ "reg"},
			Action: func (c *cli.Context) error {
				fmt.Println("selected register")
				fileReader(file)
				parseData()
				sortTransactions()
				registerCommand()
				return nil
			},
		},
		{
			Name: "print",
			Usage: "Prints out all transactions using a format readable by ledger",
			Action: func (c *cli.Context) error {
				fmt.Println("selected print")
				fileReader(file)
				parseData()
				sortTransactions()
				printCommand()
				return nil
			},
		},
	}
}

func info() {
	app.Name = "ledgerCLI"
	app.Usage = "Works for ledger's commands: balance, register, print"
	app.Authors = []*cli.Author{{Name:"David Cardenas", Email: "davidcardd1"}}
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command {{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
   `
}

func main() {
	info()
	commands()
	flags()
	sort.Sort(cli.FlagsByName(app.Flags))
    sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}