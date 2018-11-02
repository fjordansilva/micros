package main

import (
	"context"
	pb "github.com/fjordansilva/micros/user-service/proto/user"
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
	_, err := srv.repo.GetByEmailAndPassword(req)
	if err != nil {
		return err
	}

	res.Token = "testingabc"
	return nil
}

func (srv *handler) Create(ctx context.Context, req *pb.User, res *pb.Response) error  {
	if err := srv.repo.Create(req); err != nil {
		return err
	}
	res.User = req
	return nil
}

func (srv *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error  {
	return nil
}