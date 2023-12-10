package data

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	user := "root"
	password := "root123"
	hostname := "4.246.131.100"
	port := "3306"
	dbname := "empdb"

	var err error
	db, err = sql.Open("mysql", user+":"+password+"@tcp("+hostname+":"+port+")/"+dbname)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

}

func GetAllEmployees() ([]Employee, error) {
	rows, err := db.Query("SELECT empid, name, age, city FROM employee")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee

	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Age, &emp.City)
		if err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	return employees, nil
}

// GetEmployee gets a single employee by ID
func GetEmployeeByID(id int) (*Employee, error) {
	row := db.QueryRow("SELECT empid, name, age, city FROM employee WHERE empid = ?", id)

	var emp Employee
	err := row.Scan(&emp.ID, &emp.Name, &emp.Age, &emp.City)
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

// AddEmployee adds an employee to the database
func AddEmployee(emp Employee) (int64, error) {
	result, err := db.Exec("INSERT INTO employee(name, age, city) VALUES(?, ?, ?)", emp.Name, emp.Age, emp.City)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdateEmployee updates an employee in the database
func UpdateEmployee(emp Employee) (int64, error) {
	result, err := db.Exec("UPDATE employee SET name = ?, age = ?, city = ? WHERE empid = ?", emp.Name, emp.Age, emp.City, emp.ID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// DeleteEmployee deletes an employee from the database
func DeleteEmployee(id int) (int64, error) {
	result, err := db.Exec("DELETE FROM employee WHERE empid = ?", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// CloseDB closes the database connection
func CloseDB() {
	db.Close()
}
