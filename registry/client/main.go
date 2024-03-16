package main

// ----- imports
import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/oaxley/mlserver/registry/data"
	pb "github.com/oaxley/mlserver/registry/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ----- vars
var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
	port               = flag.Int("port", 50051, "The server port")
	hostname           = flag.String("hostname", "localhost", "The server address")
)

// ----- functions

// register a new service
func RegisterService(client pb.RegistryServiceClient, service *pb.ServiceDefinition) {
	log.Printf("Recording a new service to Registry server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.SetService(ctx, service)
	if err != nil {
		log.Fatalf("client: SetService failed: %v", err)
	}
	log.Println(response)
}

// query a service
func QueryService(client pb.RegistryServiceClient, query *pb.QueryService) {
	log.Printf("Querying an existing service ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.GetService(ctx, query)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Service parameters:")
	log.Println("Model Name   :", response.ModelName)
	log.Println("Model Version:", response.ModelVersion)
	log.Println("Hostname     :", response.Hostname)
	log.Println("Port         :", response.Port)
}

func main() {
	flag.Parse()

	// TLS
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// connect to gRPC server
	targetString := fmt.Sprintf("%s:%d", *hostname, *port)
	conn, err := grpc.Dial(targetString, opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// create a new client
	client := pb.NewRegistryServiceClient(conn)

	// register this new service
	service := pb.ServiceDefinition{
		ModelName:    "my_super_model",
		ModelVersion: "1.2.3",
		Hostname:     "my_server",
		Port:         12345,
	}

	RegisterService(client, &service)

	// wait 10s before requesting the client again
	time.Sleep(10 * time.Second)

	// retrieve a fake service
	query := pb.QueryService{
		ModelName:    "banana",
		ModelVersion: "1.2.3",
	}
	QueryService(client, &query)

	// wait more
	time.Sleep(5 * time.Second)

	// retrieve real model
	query = pb.QueryService{
		ModelName:    "my_super_model",
		ModelVersion: "1.2.3",
	}
	QueryService(client, &query)
}
