package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-service/proxy_server/models"
	"github.com/go-service/proxy_server/rabbit"
	"github.com/go-service/proxy_server/restapi"
	"github.com/go-service/proxy_server/restapi/operations"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

type InvoiceRequest struct {
	ClientID string `json:"client_id"`
}

func marshal(body interface{}) string {
	msg, _ := json.Marshal(body)
	return string(msg)
}

func main() {
	var portFlag = flag.Int("port", 3000, "Port to run this service on")
	var amqpUserFlag = flag.String("amqp_user", "guest", "AMQP server username")
	var amqpPassFlag = flag.String("amqp_pass", "guest", "AMQP server password")
	var amqpHostFlag = flag.String("amqp_host", "127.0.0.1", "AMQP server address")
	var amqpPortFlag = flag.Int("amqp_port", 5672, "AMQP server port")

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewRecruitProxyServerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	flag.Parse()
	server.Port = *portFlag

	// Create new instance of our queue client
	rabbitClient := rabbitclient.NewClient(*amqpUserFlag, *amqpPassFlag, *amqpHostFlag, *amqpPortFlag)

	api.GetClientsHandler = operations.GetClientsHandlerFunc(
		func(params operations.GetClientsParams) middleware.Responder {
			resp, err := rabbitClient.GetReply("backend.clients", "{}")
			if err != nil {
				return operations.NewGetClientsInternalServerError()
			}
			clients := &models.Clients{}
			err = json.Unmarshal(resp, clients)
			if err != nil {
				return operations.NewGetClientsInternalServerError()
			}
			return operations.NewGetClientsOK().WithPayload(clients)
		})
	api.GetInvoicesHandler = operations.GetInvoicesHandlerFunc(
		func(params operations.GetInvoicesParams) middleware.Responder {
			if params.ClientID != nil {
				b := marshal(InvoiceRequest{ClientID: *params.ClientID.ClientID})
				resp, err := rabbitClient.GetReply("backend.invoices", b)
				if err != nil {
					return operations.NewGetInvoicesInternalServerError()
				}
				invoices := &models.Invoices{}
				err = json.Unmarshal(resp, invoices)
				if err != nil {
					return operations.NewGetClientsInternalServerError()
				}
				return operations.NewGetInvoicesOK().WithPayload(invoices)
			}
			errMsg := marshal(ErrorMessage{Error: "You must provide client_id parameter"})
			return operations.NewGetInvoicesBadRequest().WithPayload(errMsg)
		})

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
