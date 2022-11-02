package service

import (
	"context"

	pb "github.com/Exam/customer_service/genproto/customer"
	pbp "github.com/Exam/customer_service/genproto/post"
	l "github.com/Exam/customer_service/pkg/logger"
	grpcClient "github.com/Exam/customer_service/service/grpc_client"
	"github.com/Exam/customer_service/storage"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcClient.GrpcClientI
}

func NewCustomerService(db *sqlx.DB, log l.Logger, client grpcClient.GrpcClientI) *CustomerService {
	return &CustomerService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *CustomerService) CreateCust(ctx context.Context, req *pb.CustomerReq) (*pb.CustomerResp, error) {
	customer, err := s.storage.Customer().CreateCust(req)
	if err != nil {
		s.logger.Error("error insert", l.Any("error insert customer", err))
		return &pb.CustomerResp{}, status.Error(codes.Internal, "internal error")
	}
	return customer, nil
}

func (s *CustomerService) GetCustById(ctx context.Context, req *pb.GetCustByIdReq) (*pb.GetCustomerResp, error) {
	customer, err := s.storage.Customer().GetCustById(req)
	if err != nil {
		s.logger.Error("error getting user", l.Any("error getting user", err))
		return &pb.GetCustomerResp{}, status.Error(codes.Internal, "internal error")
	}

	post, err := s.client.Post().GetCustomerPosts(ctx, &pbp.Id{Id: req.Id})
	if err != nil {
		s.logger.Error("error getting post", l.Any("error getting post info", err))
		return &pb.GetCustomerResp{}, status.Error(codes.Internal, "internal error")
	}
	for _, pst := range post.Posts {
		post.Posts = append(post.Posts, &pbp.Post{
			Id:          pst.Id,
			Name:        pst.Name,
			Description: pst.Description,
			CustomerId:  pst.CustomerId,
		})
	}
	return customer, nil
}

func (s *CustomerService) UpdateCust(ctx context.Context, req *pb.Customer) (*pb.Customer, error) {
	customer, err := s.storage.Customer().UpdateCust(req)
	if err != nil {
		s.logger.Error("error updating customer", l.Any("error updating customer info", err))
		return &pb.Customer{}, status.Error(codes.Internal, "internal error")
	}
	return customer, nil
}

func (s *CustomerService) ListCusts(context.Context, *pb.Empty) (*pb.ListCustsResp, error) {
	customer, err := s.storage.Customer().ListCusts()
	if err != nil {
		s.logger.Error("error collecting all customers", l.Any("error collecting customer info", err))
		return &pb.ListCustsResp{}, status.Error(codes.Internal, "internal error")
	}
	return customer, nil
}

func (s *CustomerService) DeleteCust(ctx context.Context, req *pb.Id) (*pb.Empty, error) {
	customer, err := s.storage.Customer().DeleteCust(req)
	if err != nil {
		s.logger.Error("error while deleting", l.Any("error while deleting customer info", err))
	}
	return customer, nil
}

func (s *CustomerService) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	boolean, err := s.storage.Customer().CheckField(req.Field, req.Value)
	if err != nil {
		s.logger.Error("error checkfield user", l.Any("error checking field", err))
		return &pb.CheckFieldResponse{}, status.Error(codes.Internal, "internal error")
	}
	return &pb.CheckFieldResponse{
		Exists: boolean.Exists,
	}, nil
}
