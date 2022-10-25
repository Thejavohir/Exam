package service

import (
	"context"

	pb "github.com/exam/review_service/genproto/review"
	l "github.com/exam/review_service/pkg/logger"
	grpcClient "github.com/exam/review_service/service/grpc_client"
	"github.com/exam/review_service/storage"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcClient.GrpcClientI
}

func NewReviewService(db *sqlx.DB, log l.Logger, client grpcClient.GrpcClientI) *ReviewService {
	return &ReviewService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.Review) (*pb.Review, error) {
	review, err := s.storage.Review().CreateReview(req)
	if err != nil {
		s.logger.Error("error insert", l.Any("error insert review", err))
		return &pb.Review{}, status.Error(codes.Internal, "internal error")
	}
	return review, nil
}

func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewReq) (*pb.Review, error) {
	review, err := s.storage.Review().GetReview(req)
	if err != nil {
		s.logger.Error("error getting review", l.Any("error getting review", err))
		return &pb.Review{}, status.Error(codes.Internal, "internal error")
	}
	return review, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.Review) (*pb.Review, error) {
	review, err := s.storage.Review().UpdateReview(req)
	if err != nil {
		s.logger.Error("error updating review", l.Any("error updating review", err))
		return &pb.Review{}, status.Error(codes.Internal, "internal error")
	}
	return &review, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	err := s.storage.Review().DeleteReview(req)
	if err != nil {
		s.logger.Error("error deleting review", l.Any("error deleting review info", err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.Empty{}, nil
}

func (s *ReviewService) GetPostReview(ctx context.Context, req *pb.GetReviewPost) (*pb.Reviews, error) {
	review, err := s.storage.Review().GetPostReview(req.PostId)
	if err != nil {
		s.logger.Error("error getting post reviews", l.Any("error getting post review info", err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return review, nil
}
