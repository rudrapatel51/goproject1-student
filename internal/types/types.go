package types

type Student struct {
	ID   int    `validate:"required" json:"id"`
	Name string `validate:"required" json:"name"`
	Age  int    `validate:"required,gte=0,lte=150" json:"age"`
}
