package repository

type Authorization interface {
}

type orderTracking interface {
}

type orderItem interface{}

type Repository struct {
	Authorization
	orderTracking
	orderItem
}

func NewRepository() *Repository {
	return &Repository{}
}
