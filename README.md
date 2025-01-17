[![GoDoc](https://godoc.org/github.com/johnbenedictyan/polluter?status.svg)](https://godoc.org/github.com/johnbenedictyan/polluter)
[![Go Report Card](https://goreportcard.com/badge/github.com/johnbenedictyan/polluter)](https://goreportcard.com/report/github.com/johnbenedictyan/polluter)
[![Build Status](https://travis-ci.org/johnbenedictyan/polluter.svg?branch=master)](https://travis-ci.org/johnbenedictyan/polluter)

# polluter

Mainly this package was created for testing purposes, to give the ability to seed a database with records from simple .yaml files. Polluter respects the order in files, so you can handle foreign_keys just by placing them in the right order.

## Usage

```go
package main

import "github.com/johnbenedictyan/polluter"

const input = `
roles:
- name: User
users:
- name: Roman
  role_id: 1
`

func TestX(t *testing.T) {
	db := prepareMySQL(t)
	defer db.Close()
	p := polluter.New(polluter.MySQLEngine(db))

	if err := p.Pollute(strings.NewReader(input)); err != nil {
		t.Fatalf("failed to pollute: %s", err)
	}

	....
}
```

## Examples

[See](https://github.com/johnbenedictyan/polluter/blob/master/polluter_test.go#L109) examples of usage with parallel testing.

## Testing

Make shure to start docker before testing.

```bash
go test
```

## Supported databases

* MySQL
* Postgres
* Redis

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!
