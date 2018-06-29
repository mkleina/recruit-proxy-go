# Recruit proxy Go service

Proxy REST API service written in Golang for [recruit-proxy](https://github.com/madsheep/recruit-proxy) service, based on RabbitMQ queue communication.

This service was built using [go-swagger](https://github.com/go-swagger/go-swagger) framework, so whole API structure is well documented using Open API 2.0 specifictaion (see: https://swagger.io/specification/).

API of this service can be easily visualized using [swagger-ui](https://github.com/swagger-api/swagger-ui) tool.

## REST API
API is by default exposed on port `3000`.

Available endpoints:
- `GET /clients.json` - List of all clients
- `POST /invoices.json` - List of all invoices for selected client, where request should be sent as JSON string `{ "client_id": "<Client ID>"}`

## Command Line Interface
Additionally, this service contains CLI tool for browsing data in human-readable format.

```shell
$ ./client -client_id google
Invoice data for google:

=== Invoice #1 ===

Total   : 200000000 USD
Services: Providing users data to the fbi
Customer: Federal Bureau of Investigation

=== Invoice #2 ===

Total   : 4000 USD
Services: Selling out users emails to ad companies
Customer: Big Bad Corporations
```

Tool is available from Docker container, just run:

```shell
sudo docker exec recruitproxygo_go-service_1 ./client
```

for available options, simply add `-h` flag.


## Installation
To install, go to root folder of this repo and run:
```
sudo docker-compose up
```
All services should start and REST API endpoints should be accessible from localhost, for example `http://127.0.0.1:3000/clients.json`

## Author
Mateusz Kleina
