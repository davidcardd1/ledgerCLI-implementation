# Ledger CLI implementation

My own Go implementation of Ledger-CLI for three basic commands: register, balance and print.
This project is being done as part of my SWE Internship at Encora.

## Documentation

### Installation
* You need to have a running version of Go in your computer
* Download or clone the code from here:
```sh
git clone https://github.com/davidcardd1/ledgerCLI-implementation
   ```
* Install the CLI in your computer:
```sh
go build .
go install .
```
* You might need to change the path of installation in your computer. First check where the CLI is installing with:
```sh
go list -f '{{.Target}}'
```
* then export that path to the executable:
```sh
export PATH=$PATH:/path/to/your/install/directory
```
* If this was the case, run (again) 
```sh
go install . 
```

* Now you can start using the app like any other CLI tool. (no need to use "./" before commands)

### How to use
* In the terminal run:
```sh
ledgerCLI-implementation -f `FILE` [flags] <command> [arguments ...]
```
* Three commands are currently working: `balance` (`bal`), `register` (`reg`) and `print`

* Three flags are currently working:
    * --sort EXPR (-S). Which can sort by: `date` (`d`), `amount` (`a`) or `payee` (`p`)
    * --file FILE (-f). Reads given ledger file. It defaults to "index.ledger"
    * --price-db FILE. Uses given file to read commodity prices


## Original work
See Ledger's full app and extensive documentation here: [Ledger Reference](https://www.ledger-cli.org/docs.html)


