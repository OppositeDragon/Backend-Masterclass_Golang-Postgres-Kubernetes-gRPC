package grpcapi

import (
	"context"
	"database/sql"
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/pb"
	util "simple_bank/util"
	"simple_bank/validator"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}
	err = util.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password: %v", err)
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %v", err)
	}
	meta := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:               refreshPayload.Id,
		Username:         user.Username,
		AccessToken:      accessToken,
		AccessExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshPayload.ExpiredAt,
		UserAgent:        util.StringPtrToSqlNullString(&meta.UserAgent),
		ClientIp:         util.StringPtrToSqlNullString(&meta.ClientIp),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	response := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}
	return response, nil
}

func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	refreshPayload, err := server.tokenMaker.VerifyToken(req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token: %v", err)
	}
	session, err := server.store.GetSession(ctx, refreshPayload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "session not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get session: %v", err)
	}

	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		return nil, status.Errorf(codes.PermissionDenied, " %v", err)
	}
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("session user mismatch")
		return nil, status.Errorf(codes.Internal, " %v", err)
	}
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("session token mismatch")
		return nil, status.Errorf(codes.Unauthenticated, " %v", err)
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}
	session, err = server.store.UpdateSessionAccess(ctx, db.UpdateSessionAccessParams{
		ID:              refreshPayload.Id,
		AccessToken:     accessToken,
		AccessExpiresAt: accessPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update session: %v", err)
	}
	response := &pb.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	}
	return response, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
