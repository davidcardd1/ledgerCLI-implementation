package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()
var file string
var comments = []rune{'!', ';', '#', '%', '|', '*'}

type Transaction struct {
	file 		string
	date 		time.Time
	payee 		string
	postings	[]*Posting

}

type Posting struct {
	transaction 	*Transaction
	account 		*Account
	commodity 		*Commodity
	commodityDate	time.Time
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
			fmt.Println(line)
		}
	}
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
			Action: func(ctx *cli.Context, s string) error {
				fileReader(file)
				return nil
			},
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