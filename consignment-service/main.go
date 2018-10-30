// consigment-service/main.go
package main

import (
	// Import the generated protobuf code
	"context"
	pb "github.com/fjordansilva/micros/consignment-service/proto/consignment"
	"github.com/micro/go-micro"
	"log"
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
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Se devuelve el mensaje de respuesta
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags
	srv.Init()

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}