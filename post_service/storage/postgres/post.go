package postgres

import (
	"fmt"
	"log"

	pb "github.com/exam/post_service/genproto/post"

	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

// NewPostRepo
func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) CreatePost(pst *pb.PostReq) (*pb.PostResp, error) {
	postResp := pb.PostResp{}
	err := r.db.QueryRow(`insert into post(
		name,
		description,
		customer_id) values($1, $2, $3) returning 
		id, 
		name, 
		description, 
		customer_id`,
		pst.Name,
		pst.Description,
		pst.CustomerId).Scan(
		&postResp.Id,
		&postResp.Name,
		&postResp.Description,
		&postResp.CustomerId,
	)
	if err != nil {
		return &pb.PostResp{}, err
	}
	for _, media := range pst.Medias {
		mediaResp := pb.Media{}
		err := r.db.QueryRow(`insert into media(
			name,
			post_id) values($1, $2) returning
			id,
			name,
			post_id`,
			media.Name,
			postResp.Id).Scan(
				&mediaResp.Id,
				&mediaResp.Name,
				&mediaResp.PostId,
			)
		if err != nil {
			fmt.Println("error while inserting media")
			return &pb.PostResp{}, err
		}
		postResp.Medias = append(postResp.Medias, &mediaResp)
	}
	return &postResp, nil
}

func (r *postRepo) UpdatePost(pst *pb.Post) (*pb.Post, error) {
	_, err := r.db.Exec(`update post SET 
	name=$1, 
	description=$2,
	cutomer_id=$3 where id = $4`,
		pst.Name,
		pst.Description,
		pst.CustomerId,
		pst.Id)
	return pst, err
}

func (r *postRepo) GetCustomerPosts(id *pb.Id) (*pb.ListPostsResp, error) {
	postResp := pb.ListPostsResp{}
	rows, err := r.db.Query(`select
		id, 
		name, 
		description, 
		customer_id from post where customer_id = $1 and deleted_at is null`, id.Id)
	if err != nil {
		log.Fatal("Error while scanning", err)
		return &pb.ListPostsResp{}, err
	}
	defer rows.Close()
	for rows.Next() {
		post := pb.Post{}
		err = rows.Scan(
			&post.Id, &post.Name, &post.Description, &post.CustomerId,
		)
		if err != nil {
			log.Fatal("Error while scanning", err)
			return &pb.ListPostsResp{}, err
		}
		postResp.Posts = append(postResp.Posts, &post)
	}
	return &postResp, nil
}

func (r *postRepo) GetPost(id *pb.GetPostReq) (*pb.GetPostResp, error) {
	postResp := pb.GetPostResp{}
	err := r.db.QueryRow(`select
		id, 
		name, 
		description, 
		customer_id from post where id = $1`, id.Id).Scan(
		&postResp.Id,
		&postResp.Name,
		&postResp.Description,
		&postResp.CustomerId,
	)
	if err != nil {
		return &pb.GetPostResp{}, err
	}

	rows, err := r.db.Query(`select id, name, post_id from media where post_id=$1`, id.Id)
	if err != nil {
		fmt.Println("error while getting media")
		return &pb.GetPostResp{}, err
	}
	defer rows.Close()

	for rows.Next() {
		media := pb.Media{}
		err = rows.Scan(
			&media.Id, 
			&media.Name, 
			&media.PostId,)
		if err != nil {
			fmt.Println("error while scanning media")
			return &pb.GetPostResp{}, err
		}
		postResp.Medias = append(postResp.Medias, &media)
	}
	fmt.Println(postResp)
	return &postResp, nil
}

func (r *postRepo) ListPosts() (*pb.ListPostsResp, error) {
	rows, err := r.db.Query(`select
		id,
		name,
		description,
		customer_id from post`)
	if err != nil {
		return &pb.ListPostsResp{}, err
	}
	defer rows.Close()

	allPosts := []*pb.Post{}
	for rows.Next() {
		allPostsResp := pb.Post{}
		err := rows.Scan(
			&allPostsResp.Id,
			&allPostsResp.Name,
			&allPostsResp.Description,
			&allPostsResp.CustomerId,
		)
		if err != nil {
			fmt.Println("error getting all cutomers", err)
			return &pb.ListPostsResp{}, err
		}
		allPosts = append(allPosts, &allPostsResp)
	}
	return &pb.ListPostsResp{Posts: allPosts}, nil
}

func (r *postRepo) DeletePost(ids *pb.Id) (*pb.Empty, error) {
	postResp := pb.Empty{}
	err := r.db.QueryRow(`update post deleted_at=NOW() where id=$1 and deleted_at is null`, ids.Id).Err()
	if err != nil {
		fmt.Println("error while deleting post")
		return &pb.Empty{}, err
	}
	return &postResp, nil
}
