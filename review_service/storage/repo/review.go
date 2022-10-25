package repo

import pb "github.com/Exam/review_service/genproto/review"

//CustomerStorageI ...
type ReviewStorageI interface {
	CreateReview(*pb.ReviewReq) (*pb.Review, error)
	GetReview(*pb.GetReviewReq) (*pb.Review, error)
	UpdateReview(*pb.Review) (pb.Review, error)
	DeleteReview(*pb.Id) error
	GetPostReview(postId int64) (*pb.Reviews, error)
}

