package postgres

import (
	"log"

	pb "github.com/Exam/review_service/genproto/review"

	"github.com/jmoiron/sqlx"
)

type reviewRepo struct {
	db *sqlx.DB
}

// NewReviewRepo
func NewReviewRepo(db *sqlx.DB) *reviewRepo {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) CreateReview(rvw *pb.Review) (*pb.Review, error) {
	revResp := pb.Review{}
	err := r.db.QueryRow(`insert into review(
		name, 
		review, 
		description, 
		post_id) values($1, $2, $3, $4) returning 
		id, 
		name, 
		review, 
		description,
		post_id`,
		rvw.Name,
		rvw.Review,
		rvw.Description,
		rvw.PostId).Scan(
		&revResp.Id,
		&revResp.Name,
		&revResp.Review,
		&revResp.Description,
		&revResp.PostId,
	)
	if err != nil {
		return &pb.Review{}, err
	}
	return &revResp, nil
}

func (r *reviewRepo) GetPostReview(postId int64) (*pb.Reviews, error) {
	reviewResp := pb.Reviews{}
	rows, err := r.db.Query(`select 
	id,
	name,
	review,
	description,
	post_id from review where post_id=$1 and deleted_at is null`, postId)
	if err != nil {
		log.Fatal("Error while scanning", err)
		return &pb.Reviews{}, err
	}
	defer rows.Close()
	for rows.Next() {
		review := pb.Review{}
		err = rows.Scan(
			&review.Id, &review.PostId, &review.Name, &review.Description, &review.Review,
		)
		if err != nil {
			log.Fatal("Error while scanning", err)
			return &pb.Reviews{}, err
		}
		reviewResp.Reviews = append(reviewResp.Reviews, &review)
	}
	return &reviewResp, nil
}

func (r *reviewRepo) GetReview(req *pb.GetReviewReq) (*pb.Review, error) {
	reviewResp := pb.Review{}
	err := r.db.QueryRow(`select 
	id,
	name,
	review,
	description,
	post_id from review where id=$1 and deleted_at is null`, req.Id).Scan(
		&reviewResp.Id,
		&reviewResp.Name,
		&reviewResp.Review,
		&reviewResp.Description,
		&reviewResp.PostId)
	if err != nil {
		log.Fatal("Error while scanning", err)
		return &pb.Review{}, err
	}
	return &reviewResp, nil
}

func (r *reviewRepo) UpdateReview(rev *pb.Review) (pb.Review, error) {
	_, err := r.db.Exec(`update review SET 
	name=$1, 
	review=$2,
	description=$3,
	post_id=$4 where id = $5 and deleted_at is null`,
		rev.Name,
		rev.Review,
		rev.Description,
		rev.PostId,
		rev.Id)
	return *rev, err
}

func (r *reviewRepo) DeleteReview(id *pb.Id) error {
	_, err := r.db.Exec(`update review set deleted_at = now() where id=$1 and deleted_at is null`, id.Id)
	return err
}
