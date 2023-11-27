package grpcapi

import (
	db "simple_bank/db/sqlc"
	"simple_bank/pb"
	util "simple_bank/util"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		Name1:             user.Name1,
		Name2:             util.SqlNullStringToString(user.Name2),
		Lastname1:         user.Lastname1,
		Lastname2:         util.SqlNullStringToString(user.Lastname2),
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
