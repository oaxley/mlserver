package main

// ----- imports
import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/oaxley/mlserver/registry/data"
	pb "github.com/oaxley/mlserver/registry/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ----- vars
var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 50051, "The server port")
	hostname = flag.String("hostname", "localhost", "The server bind address")
)

// ----- structs
type registryServiceServer struct {
	pb.UnimplementedRegistryServiceServer

	mutex    sync.Mutex // to protect access to the registry
	registry map[string]*pb.ServiceDefinition
}

// ----- functions

// register a new model
func (s *registryServiceServer) SetService(ctx context.Context, service *pb.ServiceDefinition) (*pb.Response, error) {
	// create the registry key
	key := service.ModelName + ":" + service.ModelVersion

	// insert in the registry
	s.mutex.Lock()
	s.registry[key] = service
	s.mutex.Unlock()

	// print
	log.Printf("Registered new model [%s:%s] from [%s:%d]", service.ModelName, service.ModelVersion, service.Hostname, service.Port)
	response := pb.Response{
		Message: "200 OK",
	}

	return &response, nil
}

// query for a service
func (s *registryServiceServer) GetService(cyx context.Context, service *pb.QueryService) (*pb.ServiceDefinition, error) {
	// retrieve the key
	key := service.ModelName + ":" + service.ModelVersion

	// get the value from the registry
	s.mutex.Lock()
	value, exists := s.registry[key]
	s.mutex.Unlock()

	// look if the entry exists
	if !exists {
		return nil, errors.New("entry does not exist")
	}
	return &pb.ServiceDefinition{
		ModelName:    value.ModelName,
		ModelVersion: value.ModelVersion,
		Hostname:     value.Hostname,
		Port:         value.Port,
	}, nil
}

func newServer() *registryServiceServer {
	s := &registryServiceServer{registry: make(map[string]*pb.ServiceDefinition)}
	return s
}

func main() {

	// parse the command line options
	flag.Parse()

	// start TCP server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *hostname, *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// prepare gRPC server
	var opts []grpc.ServerOption

	// TLS
	if *tls {
		if *certFile == "" {
			// default location for the TLS certificate
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			// default location for the TLS Key file
			*keyFile = data.Path("x509/server_key.pem")
		}

		// create a new credentials
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	// create new gRPC server
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterRegistryServiceServer(grpcServer, newServer())

	// start
	log.Printf("Starting registry server on %s:%d", *hostname, *port)
	grpcServer.Serve(lis)
}
