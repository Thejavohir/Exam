package grpcClient

import (
	"fmt"

	"github.com/Exam/post_service/config"
	pb "github.com/Exam/post_service/genproto/review"
	pbp "github.com/Exam/post_service/genproto/customer"
	"google.golang.org/grpc"
)

// GrpcClientI ...
type GrpcClientI interface {
	Review() pb.ReviewServiceClient
	Customer() pbp.CustomerServiceClient
}

type GrpcClient struct {
	cfg             config.Config
	reviewService   pb.ReviewServiceClient
	customerService pbp.CustomerServiceClient
}

func New(cfg config.Config) (*GrpcClient, error) {
	connCustomer, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CustomerServiceHost, cfg.CustomerServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("customer service dial host:%s, port: %d", cfg.CustomerServiceHost, cfg.CustomerServicePort)
	}
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.ReviewServiceHost, cfg.ReviewServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("review service dial host:%s, port: %d", cfg.ReviewServiceHost, cfg.ReviewServicePort)
	}
	return &GrpcClient{
		cfg:             cfg,
		reviewService:   pb.NewReviewServiceClient(connReview),
		customerService: pbp.NewCustomerServiceClient(connCustomer),
	}, nil
}

func (s *GrpcClient) Review() pb.ReviewServiceClient {
	return s.reviewService
}
func (s *GrpcClient) Customer() pbp.CustomerServiceClient {
	return s.customerService
}
