package repo

import pb "github.com/Exam/customer_service/genproto/customer"

//CustomerStorageI ...
type CustomerStorageI interface {
    CreateCust(*pb.CustomerReq) (*pb.CustomerResp, error)
	GetCustById(*pb.GetCustByIdReq) (*pb.GetCustomerResp, error)
	UpdateCust(*pb.Customer) (*pb.Customer, error)
	ListCusts() (*pb.ListCustsResp, error)
	DeleteCust(*pb.Id) (*pb.Empty, error)
	CheckField(field, value string) (*pb.CheckFieldResponse, error)  
	Login(*pb.LoginReq) (*pb.LoginResp, error)
}