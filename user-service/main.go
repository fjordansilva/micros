package main

import (
	pb "github.com/fjordansilva/micros/user-service/proto/user"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

func main() {

	// Create a database connection and handle it...
	// close before exit
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Clould not connecto to DB: %v", err)
	}

	// Automatically migrates the user struct into database columns/types, etc.
	// This will check for changes and migrate them each time this service is restarted.
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}
	tokenService := &TokenService{repo}

	// Create a new service. Include some config options
	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags
	srv.Init()

	// Register handler
	pb.RegisterUserServiceHandler(srv.Server(), &handler{repo, tokenService})

	// Run the server
	if err := srv.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
