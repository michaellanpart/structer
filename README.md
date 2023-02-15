structer
========
The `structer` facilitates the loading of environment, and default values for a struct. The loading of these values happens in the order of `environment` first then `default`, if no value is provided. Data type defaults on struct properties if not strictly set will floow the associated data type. COnfiguration of this happens in the struct tag using `default:"myvalue"` or `env:"environment property name"`.

Reasons to use this library. 
- Supports almost all kind of types
  - Scalar types
    - `int/8/16/32/64`, `uint/8/16/32/64`, `float32/64`
    - `uintptr`, `bool`, `string`
  - Complex types
    - `map`, `slice`, `struct`
  - Nested types
    - `map[K1]map[K2]Struct`, `[]map[K1]Struct[]`
  - Aliased types
    - `time.Duration`
    - e.g., `type Enum string`
  - Pointer types
    - e.g., `*SampleStruct`, `*int`
- Recursively initializes fields in a struct
- Dynamically sets default values using [`structer.Values.Set`] method
- No instantiation is required.
- Preserves non-initial values from being reset with a default value


Usage
-----

```go
// this code with associated env file are in ./example folder
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

	// format and output object to console
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

```
console output
``` json
{
    "id": 39438454,
    "name": "Roger Rabit",
    "employment_status": "Employed",
    "contacts": [
        {
            "number": "817-273-3746"
        },
        {
            "number": "415-384-9919"
        }
    ]
}

```
