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
	return CreateUser(s.DB, user)
}

func (s *Service) GetUsers() ([]User, error) {
	return GetUsers(s.DB)
}

func (s *Service) GetUserByID(id int) (*User, error) {
	return GetUserByID(s.DB, id)
}

func (s *Service) UpdateUser(id int, user *User) error {
	return UpdateUser(s.DB, id, user)
}

func (s *Service) DeleteUser(id int) error {
	return DeleteUser(s.DB, id)
}
