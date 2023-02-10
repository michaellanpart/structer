package main

import (
	"encoding/json"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/michaellanpart/structer"
)

func main() {

	// load environment variables
	godotenv.Load("./example/example.env")

	// instantiate struct
	user := &User{}

	// set env and default values
	structer.Values.Set(user)

	// output object to console
	output, _ := json.MarshalIndent(user, "", "    ")
	fmt.Printf("%s\n", output)
}

type Contact struct {
	Number string `json:"number"`
}
type User struct {
	ID               int       `default:"0" env:"EmployeeID" json:"id"`
	Name             string    `default:"" env:"EmployeeName" json:"name"`
	EmploymentStatus string    `default:"Employed" json:"employment_status"`
	Contacts         []Contact `default:"[{\"Number\":\"817-273-3746\"},{\"Number\":\"415-384-9919\"}]" json:"contacts"`
}
