// consigment-service/main.go
package main

import (
	// Import the generated protobuf code
	"context"
	pb "github.com/fjordansilva/micros/consignment-service/proto/consignment"
	vesselProto "github.com/fjordansilva/micros/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
)

type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repositorio - Dummy repository. Simula el uso de un datastore
// Se cambiara con una implementacion real
type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// El servicio debe implementar todos los metodos para satisfacer la definicion realizada en el fichero protobuf.
// Se puede comprobar la interfaz en el fichero autogenerado.
type service struct {
	repo Repository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsigment - creamos solo un metodo para el servicio
// Este metodo obtiene un contexto y una request como argumento que se envia al servidor gRPC.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	// Here we call a client instance of our vessel service with out consignment weight, and the amount of
	// containers as the capacity value
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity: int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	// We set the vesselId as the vessel we got back from our vessel service
	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
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
	repo := &ConsignmentRepository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// Init will parse the command line flags
	srv.Init()

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	// Run the server
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}