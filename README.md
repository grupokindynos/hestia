# Hestia
> Goddess of the hearth, home, domesticity, family, and the state

[![CircleCI](https://circleci.com/gh/grupokindynos/hestia.svg?style=svg)](https://circleci.com/gh/grupokindynos/hestia)
[![codecov](https://codecov.io/gh/grupokindynos/hestia/branch/master/graph/badge.svg)](https://codecov.io/gh/grupokindynos/hestia)
[![Go Report](https://goreportcard.com/badge/github.com/grupokindynos/hestia)](https://goreportcard.com/report/github.com/grupokindynos/hestia) 
[![GoDocs](https://godoc.org/github.com/grupokindynos/hestia?status.svg)](http://godoc.org/github.com/grupokindynos/hestia)

Microservice API for safe using firebase auth and mongodb
## Deploy

#### Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/grupokindynos/hestia)

#### Docker

To deploy to docker, simply pull the image
```
docker pull kindynos/hestia:latest
```
Create a new `.env` file with all the necessary environment variables defined on `app.json`

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


## API Reference
@TODO

## Testing

Simply run:
```
go test ./...
```

## Contributing

Pull requests accepted.
