package data

type Employee struct {
	ID   int    `json:"empid"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}
