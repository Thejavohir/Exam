package services

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/Exam/api_gateway/config"
	pb "github.com/Exam/api_gateway/genproto/customer"
	pbp "github.com/Exam/api_gateway/genproto/post"
	pbr "github.com/Exam/api_gateway/genproto/review"
)

type IServiceManager interface {
	CustomerService() pb.CustomerServiceClient
	PostService() pbp.PostServiceClient
	ReviewService() pbr.ReviewServiceClient
}

type serviceManager struct {
	customerService pb.CustomerServiceClient
	postService pbp.PostServiceClient
	reviewService pbr.ReviewServiceClient
}

func (s *serviceManager) CustomerService() pb.CustomerServiceClient {
	return s.customerService
}

func (s *serviceManager) PostService() pbp.PostServiceClient {
	return s.postService
}

func (s *serviceManager) ReviewService() pbr.ReviewServiceClient {
	return s.reviewService
}
																										
func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	conCustomer, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.CustomerServiceHost, conf.CustomerServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	conPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.PostServiceHost, conf.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	conReview, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.ReviewServiceHost, conf.ReviewServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		customerService: pb.NewCustomerServiceClient(conCustomer),
		postService: pbp.NewPostServiceClient(conPost),
		reviewService: pbr.NewReviewServiceClient(conReview),
	}

	return serviceManager, nil
}
