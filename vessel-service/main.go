package main

import (
	pb "github.com/fjordansilva/micros/vessel-service/proto/vessel"
)

type Repository {
	FindAvailable(*pb.Specification)
}

func main() {
	
}
