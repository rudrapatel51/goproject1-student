package storage

import (
	"errors"

	"github.com/rudrapatel51/goproject1-student/internal/types"
)

// type Student struct {
// 	ID   int    `validate:"required" json:"id"`
// 	Name string `validate:"required" json:"name"`
// 	Age  int    `validate:"required,gte=0,lte=150" json:"age"`
// }

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	Close()
}

var ErrStudentEmailAlreadyExists = errors.New("student email already exists")
var ErrStudentNotFound = errors.New("student not found")