package user

import "database/sql"

// Service предоставляет бизнес-логику для работы с пользователями
type Service struct {
	DB *sql.DB
}

// NewService создает новый экземпляр Service
func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

// ProvideUserService предоставляет новый экземпляр Service
func ProvideUserService(db *sql.DB) *Service {
	return NewService(db)
}

// Реализация методов Service
func (s *Service) CreateUser(user *User) error {
	return CreateUser(user)
}

func (s *Service) GetUsers() ([]User, error) {
	return GetUsers()
}

func (s *Service) GetUserByID(id int) (*User, error) {
	return GetUserByID(id)
}

func (s *Service) UpdateUser(id int, user *User) error {
	return UpdateUser(id, user)
}

func (s *Service) DeleteUser(id int) error {
	return DeleteUser(id)
}
