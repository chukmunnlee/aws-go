package main

import (
	"log"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_HOST string = "localhost"
	DB_PORT int = 3306
	DATABASE string = "mysql"
	DB_USER string = "fred"
	DB_PASSWORD string = "fred"
)

func CheckError(err error) {
	if nil == err {
		return
	}

	log.Fatalf("Error: %v\n", err)
}

type Customer struct {
	Id int64
	Company string
	LastName string
	FirstName string
	Email string
	JobTitle string
	BusinessPhone string
	HomePhone string
	MobilePhone string
	FaxNumber string
	Address string
	City string
	State string
	Zip string
	Country string
	WebPage string
	Notes string
	Attachments []byte
}

func ScanCustomer(row *sql.Rows) (Customer, error) {
	var (
		id int64
		company sql.NullString
		lastName sql.NullString
		firstName sql.NullString 
		email sql.NullString
		jobTitle sql.NullString
		businessPhone sql.NullString 
		homePhone sql.NullString 
		mobilePhone sql.NullString 
		faxNumber sql.NullString
		address sql.NullString
		city sql.NullString
		state_province sql.NullString
		zip sql.NullString
		country sql.NullString
		webPage sql.NullString
		notes sql.NullString
		attachments []byte
	)

	var customer Customer

	if err := row.Scan(&id, &company, &lastName, &firstName, &email, &jobTitle, &businessPhone,
			&homePhone, &mobilePhone, &faxNumber, &address, &city, &state_province, &zip,
			&country, &webPage, &notes, &attachments); nil != err {
				return customer, err
		}

		return Customer{
			Id: id,
			Company: company.String,
			LastName: lastName.String,
			FirstName: firstName.String,
			Email: email.String,
			JobTitle: jobTitle.String,
			BusinessPhone: businessPhone.String,
			HomePhone: homePhone.String,
			MobilePhone: mobilePhone.String,
			FaxNumber: faxNumber.String,
			Address: address.String,
			City: city.String,
			State: state_province.String,
			Zip: zip.String,
			Country: country.String,
			WebPage: webPage.String,
			Attachments: attachments,
		}, nil
}

func main() {

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/northwind", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT)

	log.Printf("Connecting to %s with %s\n", DATABASE, conn)

	db, err := sql.Open(DATABASE, conn)
	CheckError(err)
	defer db.Close()

	/*
	if err := db.Ping(); nil != err {
		log.Fatalf(err.Error())
	} else {
		fmt.Println("Ping successful")
	}
	*/

	//pStmt, err := db.Prepare("select * from customers limit ? offset ?")
	pStmt, err := db.Prepare("select * from customers where id = ?")
	CheckError(err)
	defer pStmt.Close()

	rows, err := pStmt.Query(12)
	CheckError(err)

	if rows.Next() {
		fmt.Println("Found record for 11")
		result, err := ScanCustomer(rows)
		CheckError(err)
		log.Printf("Customer: %v\n", result)
	} else {
		fmt.Println("No records found")
	}
}
