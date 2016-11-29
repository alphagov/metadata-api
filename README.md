# Metadata API

[![continuous integration status](https://travis-ci.org/alphagov/metadata-api.svg?branch=master)](http://travis-ci.org/alphagov/metadata-api)

A small HTTP application that acts as an easy way to get metadata
about given URLs on [GOV.UK](https://www.gov.uk/).

## Requirements

Go 1.6.2

## Development

From within `$GOPATH` run `make` to run the tests and build a binary.
Tests, build etc can also be run using the standard `go build`, `go test`...

## Dependencies

Dependencies are vendored into `vendor` and have been committed so
shouldn't need to be installed.

Use [godep](https://github.com/tools/godep) to add update
dependencies and commit to source control as they are no longer
installed during CI.

## GOV.UK Development VM

The Dev VM now has a `$GOPATH` configured at `/var/govuk/gopath` so this
app should be checked out into
`/var/govuk/gopath/src/github.com/alphagov/metadata-api` for the `go`
commands and `make` etc to function correctly.

Avoid trying to symlink the directory into the `$GOPATH` either within
the VM or on the host as this causes issues with dependencies in `vendor`. 
Working with the same directory structure on the host and allowing NFS to mount 
that into the VM works well.

For convenience using somthing like
[direnv](https://github.com/direnv/direnv) to set up a local `$GOPATH`
when you `cd` into the root of the repository keeps parity between the
host and VM.

e.g

```
export GOPATH=<path to parent gopath dir>
export PATH=$GOPATH/bin:$PATH

```
within `.envrc` in the `metadata-api` root.

## Configuration

Configuration can be handled using `ENV` variables that get
passed into the process.
