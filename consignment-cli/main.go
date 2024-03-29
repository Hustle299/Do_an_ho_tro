// consignment-cli/main.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"context"

	pb "github.com/Hustle299/Do_an_ho_tro/consignment-service/proto/consignment"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	srv := micro.NewService(micro.Name("shippy-consignment-cli"))
	srv.Init()

	client := pb.NewShippingService("shippy-consignment-service", srv.Client())

	// Contact the server and print out its response.
	file := defaultFilename
	var token string
	log.Println(os.Args)

	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments, expecting file and token")
	}

	if os.Getenv("IN_FILE") != "" {
		file = os.Getenv("IN_FILE")
	} else {
		log.Fatal("expected IN_FILE to be set")
	}

	if os.Getenv("TOKEN") != "" {
		token = os.Getenv("TOKEN")
	} else {
		log.Fatal("expected TOKEN to be set")
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	log.Printf("Attempting auth with token: %s", token)
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	r, err := client.CreateConsignment(ctx, consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	resp, err := client.GetConsignments(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not get consignments: %v", err)
	}

	for _, c := range resp.Consignments {
		log.Println(c)
	}
}
