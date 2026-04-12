package storage

import "errors"

// type Student struct {
// 	ID   int    `validate:"required" json:"id"`
// 	Name string `validate:"required" json:"name"`
// 	Age  int    `validate:"required,gte=0,lte=150" json:"age"`
// }

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
}

var ErrStudentEmailAlreadyExists = errors.New("student email already exists")