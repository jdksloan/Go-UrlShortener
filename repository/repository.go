package repository

type Repository[T any] interface {
	GetById(id int) (*T, error)
	GetByValue(val string) (*T, error)
	Insert(item *T) (*T, error)
	Update(item *T) error
	Next() (int, error)
}
