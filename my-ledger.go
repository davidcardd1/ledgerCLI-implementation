package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

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
			Action: func(ctx *cli.Context, s string) error {
				_, err := os.Open(priceDBFile)

				if err != nil {
					log.Fatal(err)
				}
				readPrices(priceDBFile)
				return nil
			},
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
				fmt.Println("Balance report")
				fileReader(file)
				parseData()
				sortTransactions()
				balanceCommand(c.Args())
				return nil
			},
		},
		{
			Name: "register",
			Usage: "Lists all postings that match the report-query with running total",
			Aliases: []string{ "reg"},
			Action: func (c *cli.Context) error {
				fmt.Println("Register report")
				fileReader(file)
				parseData()
				sortTransactions()
				registerCommand(c.Args())
				return nil
			},
		},
		{
			Name: "print",
			Usage: "Prints out all transactions using a format readable by ledger",
			Action: func (c *cli.Context) error {
				fmt.Println("Print report")
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