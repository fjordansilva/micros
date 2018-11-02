// vessel-service/main.go
package main

import (
	pb "github.com/fjordansilva/micros/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
	"os"
)

const (
	defaultHost = "localhost:27017"
)

func createDummyData(repo Repository) {
	defer repo.Close()

	vessels := []*pb.Vessel {
		&pb.Vessel{Id: "vessel01", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	// Save
	for _, v := range vessels {
		repo.Create(v)
	}
}

func main() {

	// Database host from the environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)
	// Mgo creates a 'master' session, we need to end that session
	// before the main function closes.
	defer session.Close()


	if err != nil {
		// We're wrapping the error returned from our CreateSession
		// here to add some context to the error.
		log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	}

	repo := &VesselRepository{session.Copy()}

	createDummyData(repo)

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	srv.Init()

	// Register our implementarion with
	pb.RegisterVesselServiceHandler(srv.Server(), &handler{session:session})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
