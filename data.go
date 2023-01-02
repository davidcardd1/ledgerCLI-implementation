package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	app = cli.NewApp()
	file string
	sortOption string
	priceDBFile string
	comments = []rune{'!', ';', '#', '%', '|', '*'}
	ledgerData = []string{}
	pricesData = []string{}
	transactions = []Transaction{}
	commodities = []Commodity{}
	changedCommodities = make(map[string]float64)
	root = &Account{name: "root", children: make(map[string]*Account), balance: make(map[string]float64)}
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
	name 			string
	balance			map[string]float64
	children		map[string]*Account
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
	f, err := os.Open(file)

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

func readPrices(file string) {
	f, err := os.Open(file)

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
			if strings.HasPrefix(line, "P") {
				pricesData = append(ledgerData, line)
			}
		}
	}

	for _,line := range pricesData {
		lineN := strings.Split(line, " ")
		var am float64
		am, _ = strconv.ParseFloat(lineN[4][1:], 64) 
		newCommodities(lineN[3], am)
		changedCommodities[lineN[3]] = am

	}
	changeCommodities()
}

func changeCommodities () {
	for _,t := range transactions {
		for _, p := range t.postings {
			for name, price := range changedCommodities {
				if p.commodity.name != name {
					if p.commodity.name != "$" {
						p.commodity.name = "$"
						p.amount *= price
					}
					p.commodity.name = name
					p.amount /= price
				}
			}
		}
	}
}


func newCommodities (name string, price float64) {
	if len(commodities) == 0 {
		commodities = append(commodities, Commodity{name, price})
	}
	exists := false

	for _, comms := range commodities {
		if comms.name == name {
			exists = true
		}
	}

	if !exists {
		commodities = append(commodities, Commodity{name, price})
	}
}

func (r *Account) addChildren (names []string, amount float64, comm string) {
	if len(names) == 0 {
		return
	}

	child, ok := r.children[names[0]]

	if !ok {
		child = &Account{name: names[0], children: make(map[string]*Account), balance: make(map[string]float64)}
	}

	child.balance[comm] += amount
	r.children[names[0]] = child
	
	child.addChildren(names[1:], amount, comm)
}