package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-service/proxy_client/client"
	"github.com/go-service/proxy_client/client/operations"
	"github.com/go-service/proxy_client/models"
)

func main() {
	clientID := flag.String("client_id", "", "Get invoice for specified Client ID")
	host := flag.String("host", "127.0.0.1:3000", "Recruit proxy server address")
	flag.Parse()

	// create the transport
	transport := httptransport.New(*host, "", nil)

	// create the API client, with the transport
	client := client.New(transport, strfmt.Default)

	// make the request to get all items
	if *clientID == "" {
		resp, err := client.Operations.GetClients(operations.NewGetClientsParams())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Clients data:")
		for _, client := range resp.Payload.Clients {
			fmt.Printf("- %s\n", client)
		}
	} else {
		params := operations.NewGetInvoicesParams()
		params.ClientID = &models.GetInvoicesParamsBody{ClientID: clientID}
		resp, err := client.Operations.GetInvoices(params)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Invoice data for %s:\n", *clientID)
		for i, invoice := range resp.Payload.Invoices {
			fmt.Printf("\n=== Invoice #%d ===\n\n", i+1)
			fmt.Printf("Total   : %s\n", invoice.Total)
			fmt.Printf("Services: %s\n", invoice.Services)
			fmt.Printf("Customer: %s\n", invoice.Customer)
		}
	}
}
