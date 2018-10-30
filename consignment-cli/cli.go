// consignment-cli/cli.go
package main

import (
	"context"
	"encoding/json"
	pb "github.com/fjordansilva/micros/consignment-service/proto/consignment"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"io/ioutil"
	"log"
	"os"
)

const (
	defaultFilename = "/app/consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	cmd.Init()
	// Create client
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)

	// Contact the server and print out the response
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	// Creacion de un nuevo elemento
	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", )
	}
	log.Printf("Created: %t", r.Created)

	// Obtencion de todos los elementos
	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Clould not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
