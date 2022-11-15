package postgres

import (
	"fmt"

	pb "github.com/Exam/customer_service/genproto/customer"

	"github.com/jmoiron/sqlx"
)

type customerRepo struct {
	db *sqlx.DB
}

// NewCustomerRepo
func NewCustomerRepo(db *sqlx.DB) *customerRepo {
	return &customerRepo{db: db}
}

func (r *customerRepo) CreateCust(customer *pb.CustomerReq) (*pb.CustomerResp, error) {
	customerResp := pb.CustomerResp{}
	err := r.db.QueryRow(`insert into customer(
		first_name, 
		last_name, 
		bio, 
		email, 
		phone_number,
		password) values($1, $2, $3, $4, $5, $6) returning 
		id, 
		first_name, 
		last_name, 
		bio, 
		email, 
		phone_number,
		password`,
		customer.FirstName,
		customer.LastName,
		customer.Bio,
		customer.Email,
		customer.PhoneNumber,
		customer.Password).Scan(
		&customerResp.Id,
		&customerResp.FirstName,
		&customerResp.LastName,
		&customerResp.Bio,
		&customerResp.Email,
		&customerResp.PhoneNumber,
		&customer.Password,
	)
	if err != nil {
		return &pb.CustomerResp{}, err
	}

	for _, address := range customer.Adresses {
		addressResp := pb.Address{}
		err = r.db.QueryRow(`insert into address(
			street,
			customer_id) values($1, $2) 
			returning 
			id,
			street,
			customer_id`,
			address.Street,
			customerResp.Id).Scan(
			&addressResp.Id,
			&addressResp.Street,
			&addressResp.CustomerId)
		if err != nil {
			return &pb.CustomerResp{}, err
		}
		customerResp.Adresses = append(customerResp.Adresses, &addressResp)
	}
	return &customerResp, nil
}

func (r *customerRepo) GetCustById(id *pb.GetCustByIdReq) (*pb.GetCustomerResp, error) {
	customerResp := pb.GetCustomerResp{}
	err := r.db.QueryRow(`select
		id,
		first_name, 
		last_name, 
		bio, 
		email, 
		phone_number,
		password from customer where id=$1 and deleted_at is null`, id.Id).Scan(
		&customerResp.Id,
		&customerResp.FirstName,
		&customerResp.LastName,
		&customerResp.Bio,
		&customerResp.Email,
		&customerResp.PhoneNumber,
		&customerResp.Password,
	)
	if err != nil {
		return &pb.GetCustomerResp{}, err
	}

	rows, err := r.db.Query(`select
		id,
		street,
		customer_id from address where customer_id=$1`, id.Id)
	if err != nil {
		return &pb.GetCustomerResp{}, err
	}
	defer rows.Close()

	for rows.Next() {
		address := pb.Address{}
		err = rows.Scan(
			&address.Id,
			&address.Street,
			&address.CustomerId,
		)
		if err != nil {
			return &pb.GetCustomerResp{}, err
		}
		customerResp.Adresses = append(customerResp.Adresses, &address)
	}
	return &customerResp, nil
}

func (r *customerRepo) UpdateCust(cust *pb.Customer) (*pb.Customer, error) {
	_, err := r.db.Exec(`update customer SET 
	first_name=$1, 
	last_name=$2,
	bio=$3,
	email=$4, 
	phone_number=$5,
	password=$6 where id = $7`,
		cust.FirstName,
		cust.LastName,
		cust.Bio,
		cust.Email,
		cust.PhoneNumber,
		cust.Password,
		cust.Id)
	return cust, err
}

func (r *customerRepo) ListCusts() (*pb.ListCustsResp, error) {
	rows, err := r.db.Query(`select
		id,
		first_name,
		last_name,
		bio,
		email,
		phone_number,
		password from customer`)

	if err != nil {
		return &pb.ListCustsResp{}, err
	}
	defer rows.Close()

	allCustomers := []*pb.Customer{}
	for rows.Next() {
		allCustomersResp := pb.Customer{}
		err := rows.Scan(
			&allCustomersResp.Id,
			&allCustomersResp.FirstName,
			&allCustomersResp.LastName,
			&allCustomersResp.Bio,
			&allCustomersResp.Email,
			&allCustomersResp.PhoneNumber,
		)
		if err != nil {
			fmt.Println("error getting all cutomers", err)
			return &pb.ListCustsResp{}, err
		}
		allCustomers = append(allCustomers, &allCustomersResp)
	}
	return &pb.ListCustsResp{Customers: allCustomers}, nil
}

func (r *customerRepo) DeleteCust(ids *pb.Id) (*pb.Empty, error) {
	custResp := pb.Empty{}
	err := r.db.QueryRow(`update customer set deleted_at=NOW() where id=$1 and deleted_at is null`, ids.Id).Err()
	if err != nil {
		fmt.Println("error while deleting customer")
		return &pb.Empty{}, err
	}
	return &custResp, nil
}

func (r *customerRepo) CheckField(field, value string) (*pb.CheckFieldResponse, error) {
	query := fmt.Sprintf("select count(1) from customer where %s = $1", field)
	var exists int
	err := r.db.QueryRow(query, value).Scan(&exists)
	if err != nil {
		return &pb.CheckFieldResponse{}, err
	}
	if exists == 0 {
		return &pb.CheckFieldResponse{Exists: false}, nil
	}
	return &pb.CheckFieldResponse{Exists: true}, nil
}

func (r *customerRepo) Login(req *pb.LoginReq) (*pb.LoginResp, error) {
	resp := &pb.LoginResp{}
	err := r.db.QueryRow(`select
	id,
	first_name,
	last_name,
	bio,
	password,
	phone_number from customer where email=$1`, req.Email).Scan(
		&resp.Id,
		&resp.FirstName,
		&resp.LastName,
		&resp.Bio,
		&resp.Password,
		&resp.PhoneNumber,
		&resp.RefreshToken,
	)
	if err != nil {
		fmt.Println("error while getting user login")
		return &pb.LoginResp{}, err
	}

	rows, err := r.db.Query(`SELECT id, customer_id, street FROM address WHERE customer_id=$1`, resp.Id)
	if err != nil {
		fmt.Println("error while getting addresses login")
	}
	for rows.Next() {
		address := pb.Address{}
		err = rows.Scan(&address.Id, &address.CustomerId, &address.Street)
		if err != nil {
			fmt.Println("error while scanning address")
			return &pb.LoginResp{}, err
		}
		resp.Addresses = append(resp.Addresses, &address)
	}
	return &pb.LoginResp{}, nil

}
