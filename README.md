# Metadata API

[![continuous integration status](https://travis-ci.org/alphagov/metadata_api.svg?branch=master)](http://travis-ci.org/alphagov/metadata_api)

A small HTTP application that acts as an easy way to get metadata
about given URLs on [GOV.UK](https://www.gov.uk/).

## Requirements

To run the code you will need to have at least Go 1.2 installed.

## Development

You can run the tests locally by running the following to use `godep`
to fetch the dependencies:

```bash
make
```

Alternatively you can install the dependencies directly to your
`$GOPATH` and run the tests using:

```bash
go get -v ./...
go test -v ./...
```

## Running

You can build a binary using either `make build` or `go build`. You
should then be able to run the binary directly.
