# Hestia
> Goddess of the hearth, home, domesticity, family, and the state

![Actions](https://github.com/grupokindynos/hestia/workflows/Hestia/badge.svg)
[![codecov](https://codecov.io/gh/grupokindynos/hestia/branch/master/graph/badge.svg)](https://codecov.io/gh/grupokindynos/hestia)
[![Go Report](https://goreportcard.com/badge/github.com/grupokindynos/hestia)](https://goreportcard.com/report/github.com/grupokindynos/hestia) 
[![GoDocs](https://godoc.org/github.com/grupokindynos/hestia?status.svg)](http://godoc.org/github.com/grupokindynos/hestia)

Hestia is a microservice API for safe using firebase auth and firestore

## Deploy

#### Docker

To deploy to docker, simply pull the image
```
docker pull kindynos/hestia:latest
```
Run the docker image
```
docker run -p 8080:8080 --env-file .env kindynos/hestia:latest 
```

## Building

To run Hestia from the source code, first you need to install golang, follow this guide:
```
https://golang.org/doc/install
```

To run Hestia simply clone de repository:

```
git clone https://github.com/grupokindynos/hestia 
```

Install dependencies
```
go mod download
```

Make sure the port is configured under en enviroment variable `PORT=8080`

## Running flags
```
-local
```

Set this flag to change from production database to testing database.

## API Reference

Documentation: [API Reference](https://documenter.getpostman.com/view/4345063/SVfUsSTD?version=latest)

## Testing

Run:
```
go test ./ -coverprofile=coverage.txt -covermode=atomic
go test ./models -coverprofile=models/coverage.txt -covermode=atomic
go test ./controllers -coverprofile=controllers/coverage.txt -covermode=atomic
sed '1d' ./models/coverage.txt >> ./coverage.txt &&
sed '1d' ./controllers/coverage.txt >> ./coverage.txt
```

## Contributing

To contribute to this repository, please fork it, create a new branch and submit a pull request.
