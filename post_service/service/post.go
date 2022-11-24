package service

import (
	"context"

	pb "github.com/Exam/post_service/genproto/post"
	rs "github.com/Exam/post_service/genproto/review"
	l "github.com/Exam/post_service/pkg/logger"
	grpcClient "github.com/Exam/post_service/service/grpc_client"
	"github.com/Exam/post_service/storage"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcClient.GrpcClientI
}

func NewPostService(db *sqlx.DB, log l.Logger, client grpcClient.GrpcClientI) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.PostReq) (*pb.PostResp, error) {
	post, err := s.storage.Post().CreatePost(req)
	if err != nil {
		s.logger.Error("error insert", l.Any("error insert post", err))
		return &pb.PostResp{}, status.Error(codes.Internal, "internal error")
	}
	return post, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	post, err := s.storage.Post().UpdatePost(req)
	if err != nil {
		s.logger.Error("error updating post", l.Any("error updating post", err))
		return &pb.Post{}, status.Error(codes.Internal, "internal error")
	}
	return post, nil
}

func (s *PostService) GetPost(ctx context.Context, req *pb.GetPostReq) (*pb.GetPostResp, error) {
	post, err := s.storage.Post().GetPost(req)
	if err != nil {
		s.logger.Error("error getting post", l.Any("error getting post info", err))
		return &pb.GetPostResp{}, status.Error(codes.Internal, "internal error")
	}
	rev, err := s.client.Review().GetPostReview(ctx, &rs.GetReviewPost{PostId: req.Id})

	if err != nil {
		s.logger.Error("error getting post", l.Any("error getting post info", err))
		return &pb.GetPostResp{}, status.Error(codes.Internal, "internal error")
	}
	for _, r := range rev.Reviews {
		rev.Reviews = append(rev.Reviews, &rs.Review{
			Id:          r.Id,
			PostId:      r.PostId,
			Name:        r.Name,
			Review:      r.Review,
			Description: r.Description,
		})
	}

	return post, nil

}

func (s *PostService) GetCustomerPosts(ctx context.Context, req *pb.Id) (*pb.ListPostsResp, error) {
	post, err := s.storage.Post().GetCustomerPosts(req)
	if err != nil {
		s.logger.Error("error getting customer posts", l.Any("error getting customer posts", err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return post, nil
}

func (s *PostService) ListPosts(context.Context, *pb.Empty) (*pb.ListPostsResp, error) {
	post, err := s.storage.Post().ListPosts()
	if err != nil {
		s.logger.Error("error collecting posts", l.Any("error collecting post info", err))
		return &pb.ListPostsResp{}, status.Error(codes.Internal, "internal error")
	}
	return post, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	post, err := s.storage.Post().DeletePost(req)
	if err != nil {
		s.logger.Error("error while deleting", l.Any("error while deleting post info", err))
	}
	return post, nil
}
