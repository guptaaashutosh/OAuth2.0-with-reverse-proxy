package model

// type User struct {
// 	Id int64
// 	Name string
// 	Age int64
// 	Address string
// }

// type Login struct {
// 	Id int64
// 	Name string
// }

type User struct {
	Id       int64 		`json:id`
	Email    string		`json:email`
	Password string		`json:password`
	Name     string		`json:name`
	Age      string		`json:age`
	Address  string		`json:address`
	Sid      []int64	`json:sid`
}
type GetUser struct {
	Id      int64
	Email   string
	Name    string    `json:",omitempty"` //Go's encoding/json package to control how fields are encoded in JSON
	Age     int64    `json:",omitempty"`
	Address string    `json:",omitempty"`
	Service []Service `json:",omitempty"`
}

type Login struct {
	Email    string
	Password string
}

type Service struct {
	Sid     int64
	Service string
}

