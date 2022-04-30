package user

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
	Authenticate(string, string) (User, error)
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

func (s userService) Authenticate(email string, password string) (user User, err error) {
	u, err := s.GetByEmail(email)
	if err != nil {
		return
	}
	if u.Password == "" {
		err = fmt.Errorf("user has no password set")
		return
	}
	if password == "" {
		err = fmt.Errorf("password is empty")
		return
	}
	hashed, compare := []byte(u.Password), []byte(password)
	if err = bcrypt.CompareHashAndPassword(hashed, compare); err != nil {
		return
	}
	return u, nil
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

	if user.Password != "" {
		password := []byte(user.Password)
		hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
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

	// If the password is empty, leave the current password alone
	if user.Password == "" {
		user.Password = check.Password
	}

	// A new, non-hashed password will give an error during the cost
	// calculation. As long as the password is non-empty, hash the new password
	// and store the updated value
	password := []byte(user.Password)
	if _, err := bcrypt.Cost(password); len(password) > 0 && err != nil {
		hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
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
