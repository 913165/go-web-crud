package main

import (
	"GO_CRUD_EMPLOYEES/data"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	// Initialize the database connection
	data.InitDB()
	defer data.CloseDB()

	// Endpoint to get all employees or create a new employee
	http.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Get all employees
			employees, err := data.GetAllEmployees()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(employees)
		case http.MethodPost:
			// Create new employee
			// Implement code to create a new employee using data.CreateEmployee()
			var newEmployee data.Employee
			err := json.NewDecoder(r.Body).Decode(&newEmployee)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			insertedID, err := data.AddEmployee(newEmployee)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Retrieve the newly added employee by ID for response
			createdEmployee, err := data.GetEmployeeByID(int(insertedID))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(createdEmployee)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	})

	// Endpoint to get, update, or delete an employee by ID
	http.HandleFunc("/employees/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			urlParts := strings.Split(r.URL.Path, "/")
			if len(urlParts) != 3 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid URL"))
				return
			}
			employeeID, err := strconv.Atoi(urlParts[2])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid employee ID"))
				return
			}

			switch r.Method {
			case http.MethodGet:
				// Get employee by ID
				employee, err := data.GetEmployeeByID(employeeID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(employee)

			case http.MethodPut:
				// Update employee by ID
				var updatedEmployee data.Employee
				err := json.NewDecoder(r.Body).Decode(&updatedEmployee)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				updatedEmployee.ID = employeeID // Set the ID to the specified employee ID
				rowsAffected, err := data.UpdateEmployee(updatedEmployee)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if rowsAffected == 0 {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Employee not found"))
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Employee updated successfully"))

			case http.MethodDelete:
				// Delete employee by ID
				rowsAffected, err := data.DeleteEmployee(employeeID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if rowsAffected == 0 {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Employee not found"))
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Employee deleted successfully"))
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	})

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
