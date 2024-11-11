package user

import "database/sql"

// Service — це шар логіки для роботи з користувачами
type Service struct {
	DB *sql.DB
}

// NewService створює новий екземпляр Service
func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func ProvideUserService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func ProvideUserHandler(service *Service) *Handler {
	return &Handler{Service: service}
}
