package postgres

import (
	"fmt"

	pb "github.com/exam/customer_service/genproto/customer"

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
		phone_number) values($1, $2, $3, $4, $5) returning 
		id, 
		first_name, 
		last_name, 
		bio, 
		email, 
		phone_number`,
		customer.FirstName,
		customer.LastName,
		customer.Bio,
		customer.Email,
		customer.PhoneNumber).Scan(
		&customerResp.Id,
		&customerResp.FirstName,
		&customerResp.LastName,
		&customerResp.Bio,
		&customerResp.Email,
		&customerResp.PhoneNumber,
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
		phone_number from customer where id=$1 and deleted_at is null`, id.Id).Scan(
		&customerResp.Id,
		&customerResp.FirstName,
		&customerResp.LastName,
		&customerResp.Bio,
		&customerResp.Email,
		&customerResp.PhoneNumber,
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
		customerResp.Adresses=append(customerResp.Adresses, &address)
	}
	return &customerResp, nil
}

func (r *customerRepo) UpdateCust(cust *pb.Customer) (*pb.Customer, error) {
	_, err := r.db.Exec(`update customer SET 
	first_name=$1, 
	last_name=$2,
	bio=$3,
	email=$4, 
	phone_number=$5 where id = $6`,
		cust.FirstName,
		cust.LastName,
		cust.Bio,
		cust.Email,
		cust.PhoneNumber,
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
		phone_number from customer`)

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
