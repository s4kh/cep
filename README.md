# cep
Cron expression parser

## Prerequisites
Install golang on you machine from https://golang.org/dl/. You can follow the instructions [here](https://golang.org/doc/install) based on your OS.

## To run
1. Go into projects directory.
2. Install dependencies using `go get`.
3. Build using `go build`.
4. Run example: `./cep "*/15 0 1,5 * 1-5 /usr/bin/find"`

Note: Don't forget to put quotes for cron expression arguments.

# To run tests
`go test -v`