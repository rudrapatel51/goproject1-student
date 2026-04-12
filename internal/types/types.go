package types

type Student struct {
	ID    int    `json:"id"`
	Name  string `validate:"required" json:"name"`
	Email string `validate:"required,email" json:"email"`
	Age   int    `validate:"required,gte=0,lte=150" json:"age"`
}
