package repo

import pb "github.com/exam/post_service/genproto/post"

//CustomerStorageI ...
type PostStorageI interface {
    CreatePost(*pb.PostReq) (*pb.PostResp, error)
	UpdatePost(*pb.Post) (*pb.Post, error)
	GetPost(*pb.GetPostReq) (*pb.GetPostResp, error)
	ListPosts() (*pb.ListPostsResp, error)
	DeletePost(*pb.Id) (*pb.Empty, error)
	GetCustomerPosts(*pb.Id) (*pb.ListPostsResp, error)
}