package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	app = cli.NewApp()
	file string
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

func parseData () {

	transaction := Transaction{}

	for _,line := range ledgerData {
		transInfo := strings.Split(line, " ")
	
		if matches, _ := regexp.MatchString(`\d{4}\/(1[0-2]|[1-9])\/(3[0-1]|[1-2][0-9]|[1-9])$`, transInfo[0]); matches {
			const layout = "2006/1/2"
			date, _ := time.Parse(layout, transInfo[0])
			payee := strings.Join(transInfo[1:]," ")
			
			//fmt.Printf("%v %v \n", date, payee)

			transactions = append(transactions, transaction)
			transaction = Transaction{}
			transaction.date = date
			transaction.payee = payee

		} else {
			postingInfo := strings.Split(line, "\t")

			account := postingInfo[1]
			quantity := postingInfo[2:]
			var commodity string
			var amount float64

			for _,x := range quantity {
				if x != "" {
					if x[0] == '$' || strings.HasPrefix(x, "-$"){
						commodity = "$"
						amount, _ = strconv.ParseFloat(strings.Replace(x, "$", "", 1), 64)
					} else {
						s := strings.Split(x, " ")
						commodity = s[1]
						amount, _= strconv.ParseFloat(s[0], 64)
					}
					break
				}
			}

			// fmt.Printf("%q \n", account)
			// fmt.Printf("%q \n", commodity)
			// fmt.Printf("%v \n", amount)

			newAccount := Account{account, true}
			newCommodity := Commodity{name: commodity}

			accounts = append(accounts, newAccount)
			commodities = append(commodities, newCommodity)

			newPosting := Posting{newAccount, newCommodity, amount}

			transaction.postings = append(transaction.postings, newPosting)
			newPosting = Posting{}
		}
	}
	transactions = append(transactions, transaction)
}

func flags() {
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "sort",
			Aliases: []string{"s", "S"},
			Value:   "date",
			Usage:   "Sort report using the value expression: `VEXPR`",
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

				fmt.Printf("%v", transactions)
				return nil
			},
		},
	}
}

func info() {
	app.Name = "ledgerCLI"
	app.Usage = "Works for ledger's commands: balance, register, print"
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