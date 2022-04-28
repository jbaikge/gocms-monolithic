package user

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	DisplayName string             `json:"display_name" bson:"display_name" form:"display_name"`
	Email       string             `json:"email" bson:"email" form:"email"`
	Password    string             `json:"password" bson:"password" form:"password"`
	Active      bool               `json:"active" bson:"active" form:"active"`
}

type UserRepository interface {
	GetUserByEmail(string) (User, error)
	GetUserById(primitive.ObjectID) (User, error)
	InsertUser(*User) error
	UpdateUser(*User) error
}

type UserService interface {
	GetByEmail(string) (User, error)
	GetById(primitive.ObjectID) (User, error)
	Insert(*User) error
	Update(*User) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return userService{
		repo: repo,
	}
}

func (s userService) GetByEmail(email string) (User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s userService) GetById(id primitive.ObjectID) (User, error) {
	return s.repo.GetUserById(id)
}

func (s userService) Insert(user *User) (err error) {
	if err = s.Validate(user); err != nil {
		return
	}

	check, _ := s.GetByEmail(user.Email)
	if !check.Id.IsZero() {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	if !user.Id.IsZero() {
		return fmt.Errorf("user already has an ID")
	}

	return s.repo.InsertUser(user)
}

func (s userService) Update(user *User) (err error) {
	if err = s.Validate(user); err != nil {
		return
	}

	if user.Id.IsZero() {
		return fmt.Errorf("user has no ID")
	}

	check, _ := s.GetByEmail(user.Email)
	if !check.Id.IsZero() && check.Id != user.Id {
		return fmt.Errorf("email already used by another user: %s", user.Email)
	}

	return s.repo.UpdateUser(user)
}

func (s userService) Validate(user *User) (err error) {
	if user.Email == "" {
		return fmt.Errorf("email is empty")
	}

	if user.DisplayName == "" {
		return fmt.Errorf("display name is empty")
	}

	return
}
