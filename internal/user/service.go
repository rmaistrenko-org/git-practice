package user

import "database/sql"

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func ProvideUserService(db *sql.DB) *Service {
	return NewService(db)
}

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
