[![Build Status](https://travis-ci.org/dotchev/sm.svg?branch=master)](https://travis-ci.org/dotchev/sm)

# sm

## Development

Run Postgres locally in Docker
```sh
docker run --name postgres -p 5432:5432 -d postgres
```
Install dependencies
```sh
dep ensure
```
Run the tests
```sh
go test ./... -v
```
