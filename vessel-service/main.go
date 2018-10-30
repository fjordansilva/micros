// vessel-service/main.go
package main

import (
	"context"
	pb "github.com/fjordansilva/micros/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"log"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

// FindAvailable - comprueba una especificacion contra un mapa de Vessels
// Si la capacidad y el peso maximo estan por debajo de la capacidad del vessel y su peso maximo,
// se devuelve ese vessel
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error)  {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New(1,"No vessel found by that spec", 1)
}

// gRPC service handler
type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error  {
	// Find the next available vessel
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	// Set the vessel as part of the response message type
	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel {
		&pb.Vessel{Id: "vessel01", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels:vessels}

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	srv.Init()

	// Register our implementarion with
	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
