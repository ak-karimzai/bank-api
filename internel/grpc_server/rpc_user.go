package grpcserver

import (
	"context"
	"time"

	"github.com/ak-karimzai/bank-api/internel/db"
	errorhandler "github.com/ak-karimzai/bank-api/internel/error_handler"
	"github.com/ak-karimzai/bank-api/internel/pb"
	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/ak-karimzai/bank-api/internel/validators"
	"github.com/ak-karimzai/bank-api/internel/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *GRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRespone, error) {
	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPwd, err := util.HashPasswrod(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "failed to hash pwd: %v", err)
	}

	arg := db.CreateUserParams{
		Username:  req.GetUsername(),
		HashedPwd: hashedPwd,
		FullName:  req.GetFullName(),
		Email:     req.GetEmail(),
	}

	user, err := server.UserRepo.CreateUser(ctx, arg)
	if err != nil {
		finalErr := errorhandler.DbErrorHandler(err)
		return nil, status.Errorf(toGrpcError(finalErr), finalErr.Message)
	}

	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = server.TaskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "failed to distribute task to send verify email: %v", err)
	}
	response := &pb.CreateUserRespone{
		User: converUser(user),
	}

	return response, nil
}

func (server *GRPCServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserRespone, error) {
	if violations := validateLoginUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.UserRepo.GetUser(ctx, req.GetUsername())
	if err != nil {
		finalErr := errorhandler.DbErrorHandler(err)
		return nil, status.Errorf(toGrpcError(finalErr), finalErr.Message)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPwd), []byte(req.GetPassword())); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorect password!")
	}

	accessToken, accessPayload, err := server.TokenMaker.CreateToken(
		req.GetUsername(), server.Config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token")
	}

	refreshToken, refreshPayload, err := server.TokenMaker.CreateToken(
		req.GetUsername(), server.Config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to refresh token")
	}

	mtdt := server.extractMetadata(ctx)

	session, err := server.SessionRepo.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rsp := &pb.LoginUserRespone{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  converUser(user),
	}
	return rsp, nil
}

func (server *GRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserRespone, error) {
	payload, err := server.authUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateUpdateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if payload.Username != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "cannot update other user's info")
	}

	var hashedPwd pgtype.Text
	var pwdChangedAt pgtype.Timestamptz
	if req.Password != nil {
		hPwd, err := util.HashPasswrod(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(
				codes.Internal, "failed to hash pwd: %v", err)
		}
		hashedPwd = pgtype.Text{
			String: hPwd,
			Valid:  true,
		}
		pwdChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}
	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		HashedPwd:    hashedPwd,
		PwdChangedAt: pwdChangedAt,
	}

	user, err := server.UserRepo.UpdateUser(ctx, arg)
	if err != nil {
		finalErr := errorhandler.DbErrorHandler(err)
		return nil, status.Errorf(toGrpcError(finalErr), finalErr.Message)
	}

	response := &pb.UpdateUserRespone{
		User: converUser(user),
	}

	return response, nil
}

func validateCreateUserRequest(
	req *pb.CreateUserRequest) (
	violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, filedViolation("username", err))
	}

	if err := validators.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, filedViolation("full_name", err))
	}

	if err := validators.ValidatePwd(req.GetPassword()); err != nil {
		violations = append(violations, filedViolation("password", err))
	}

	if err := validators.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, filedViolation("email", err))
	}
	return
}

func validateUpdateUserRequest(
	req *pb.UpdateUserRequest) (
	violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, filedViolation("username", err))
	}

	if req.FullName != nil {
		if err := validators.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, filedViolation("full_name", err))
		}
	}

	if req.Password != nil {
		if err := validators.ValidatePwd(req.GetPassword()); err != nil {
			violations = append(violations, filedViolation("password", err))
		}
	}

	if req.Email != nil {
		if err := validators.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, filedViolation("email", err))
		}
	}
	return
}

func validateLoginUserRequest(
	req *pb.LoginUserRequest) (
	violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, filedViolation("username", err))
	}

	if err := validators.ValidatePwd(req.GetPassword()); err != nil {
		violations = append(violations, filedViolation("password", err))
	}
	return
}
