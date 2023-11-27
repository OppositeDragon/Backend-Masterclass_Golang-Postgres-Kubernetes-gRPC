package grpcapi

import (
	"context"
	db "simple_bank/db/sqlc"
	"simple_bank/pb"
	util "simple_bank/util"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		Name1:          req.GetName1(),
		Name2:          util.StringToSqlNullString(req.GetName2()),
		Lastname1:      req.GetLastname1(),
		Lastname2:      util.StringToSqlNullString(req.GetLastname2()),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, value := err.(*pq.Error); value {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists")
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "email already in use")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}
	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}
