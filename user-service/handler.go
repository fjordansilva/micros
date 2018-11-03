package main

import (
	"context"
	pb "github.com/fjordansilva/micros/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type handler struct {
	repo Repository
	tokenService Authable
}

func (srv *handler) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := srv.repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (srv *handler) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := srv.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (srv *handler) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.repo.GetByEmailAndPassword(req)
	log.Println(user)
	if err != nil {
		return err
	}

	// Compares our given password against the hashed password stores in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

func (srv *handler) Create(ctx context.Context, req *pb.User, res *pb.Response) error  {

	// Generates a hashed version of our password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashedPass)
	if err := srv.repo.Create(req); err != nil {
		return err
	}

	token, err := srv.tokenService.Encode(req)
	if err != nil {
		return err
	}
	res.User = req
	res.Token = &pb.Token{Token: token}
	return nil
}

func (srv *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error  {
	return nil
}