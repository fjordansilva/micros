// consigment-service/main.go
package main

import (
	// Import the generated protobuf code
	"context"
	pb "github.com/fjordansilva/micros/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repositorio - Dummy repository. Simula el uso de un datastore
// Se cambiara con una implementacion real
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// El servicio debe implementar todos los metodos para satisfacer la definicion realizada en el fichero protobuf.
// Se puede comprobar la interfaz en el fichero autogenerado.
type service struct {
	repo IRepository
}

// CreateConsigment - creamos solo un metodo para el servicio
// Este metodo obtiene un contexto y una request como argumento que se envia al servidor gRPC.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Se devuelve el mensaje de respuesta
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments: consignments}, nil
}

func main() {
	repo := &Repository{}

	// Setup el servidor gRPC
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to lister: %v", err)
	}
	s := grpc.NewServer()

	// Registro de nuestro servicio en el servidor gRPC.
	// De esta forma se une la implementación al codigo autogenerado
	pb.RegisterShippingServiceServer(s, &service{repo})

	// Registro del sistema de reflexion del servidor gRPC
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}