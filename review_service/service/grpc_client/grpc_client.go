package grpcClient

import (
	"fmt"

	"github.com/exam/review_service/config"
	pbp "github.com/exam/review_service/genproto/post"
	"google.golang.org/grpc"
)

// GrpcClientI ...
type GrpcClientI interface {
	Post() pbp.PostServiceClient
}

// GrpcClient ...
type GrpcClient struct {
	cfg         config.Config
	postService pbp.PostServiceClient
}

// New ...
func New(cfg config.Config) (*GrpcClient, error) {
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("post service dial host:%s, port: %d", cfg.PostServiceHost, cfg.PostServicePort)
	}
	return &GrpcClient{
		cfg:         cfg,
		postService: pbp.NewPostServiceClient(connPost),
	}, nil
}

func (s *GrpcClient) Post() pbp.PostServiceClient {
	return s.postService
}
