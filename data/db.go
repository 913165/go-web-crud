package data

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitDB() {

	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "root123")
	hostname := getEnv("DB_HOST", "172.178.82.33")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "empdb")

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
