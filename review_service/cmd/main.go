package main

import (
	"net"

	pb "github.com/Exam/review_service/genproto/review"

	"github.com/Exam/review_service/config"
	"github.com/Exam/review_service/pkg/db"
	"github.com/Exam/review_service/pkg/logger"
	"github.com/Exam/review_service/service"
	grpcClient "github.com/Exam/review_service/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "review_service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	conDb, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx failed to connect to the database", logger.Error(err))
	}

	grpcClient, err := grpcClient.New(cfg)
	if err != nil {
		log.Fatal("error while grpc connection with client", logger.Error(err))
	}

	reviewService := service.NewReviewService(conDb, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterReviewServiceServer(s, reviewService)
	log.Info("main: server running", logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("error while listening %v", logger.Error(err))
	}
}