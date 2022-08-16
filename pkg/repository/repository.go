package repository

type ZodiacAPI interface {
}

type Repository struct {
	ZodiacAPI
}

func NewRepository() *Repository {
	return &Repository{}
}
