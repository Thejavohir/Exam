package grpcClient

import (
	"fmt"

	"github.com/Exam/customer_service/config"
	pb "github.com/Exam/customer_service/genproto/post"
	pbp "github.com/Exam/customer_service/genproto/review"

	"google.golang.org/grpc"
)

// GrpcClientI ...
type GrpcClientI interface {
	Post() pb.PostServiceClient
	Review() pbp.ReviewServiceClient
}

// GrpcClient ...
type GrpcClient struct {
	cfg           config.Config
	postService   pb.PostServiceClient
	reviewService pbp.ReviewServiceClient
}

// New ...
func New(cfg config.Config) (*GrpcClient, error) {
	conPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("post service dial host: %s port: %d", cfg.PostServiceHost, cfg.PostServicePort)
	}

	conReview, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.ReviewServiceHost, cfg.ReviewServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("review service dial host:%s, port: %d", cfg.ReviewServiceHost, cfg.ReviewServicePort)
	}

	return &GrpcClient{
		cfg:           cfg,
		postService:   pb.NewPostServiceClient(conPost),
		reviewService: pbp.NewReviewServiceClient(conReview),
	}, nil
}

func (r *GrpcClient) Post() pb.PostServiceClient {
	return r.postService
}

func (s *GrpcClient) Review() pbp.ReviewServiceClient {
	return s.reviewService
}
