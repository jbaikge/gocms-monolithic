package user

import (
	"fmt"
	"testing"

	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ UserRepository = mockUserRepository{}

type mockUserRepository struct {
	byEmail map[string]User
	byId    map[primitive.ObjectID]User
}

func NewMockUserRepository() mockUserRepository {
	return mockUserRepository{
		byEmail: make(map[string]User),
		byId:    make(map[primitive.ObjectID]User),
	}
}

func (r mockUserRepository) DeleteUser(id primitive.ObjectID) (err error) {
	user, ok := r.byId[id]
	if !ok {
		// Silent failure
		return
	}
	delete(r.byId, id)
	delete(r.byEmail, user.Email)
	return
}

func (r mockUserRepository) GetUserByEmail(email string) (user User, err error) {
	user, ok := r.byEmail[email]
	if !ok {
		err = fmt.Errorf("user not found: %s", email)
	}
	return
}

func (r mockUserRepository) GetUserById(id primitive.ObjectID) (user User, err error) {
	user, ok := r.byId[id]
	if !ok {
		err = fmt.Errorf("user not found: %s", id)
	}
	return
}

func (r mockUserRepository) InsertUser(user *User) (err error) {
	user.Id = primitive.NewObjectID()
	r.byId[user.Id] = *user
	r.byEmail[user.Email] = *user
	return
}

func (r mockUserRepository) UpdateUser(user *User) (err error) {
	if err = r.DeleteUser(user.Id); err != nil {
		return
	}
	r.byId[user.Id] = *user
	r.byEmail[user.Email] = *user
	return
}

func TestGetByEmail(t *testing.T) {
	service := NewUserService(NewMockUserRepository())

	user := User{DisplayName: "Test Testerly", Email: "test@test.com"}
	assert.NoError(t, service.Insert(&user))

	check, err := service.GetByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.Id, check.Id)

	_, err = service.GetByEmail("moo@cow.com")
	assert.Error(t, err)
}

func TestGetById(t *testing.T) {
	service := NewUserService(NewMockUserRepository())

	user := User{DisplayName: "Test Testerly", Email: "test@test.com"}
	assert.NoError(t, service.Insert(&user))

	check, err := service.GetById(user.Id)
	assert.NoError(t, err)
	assert.Equal(t, user.Id, check.Id)

	_, err = service.GetById(primitive.NewObjectID())
	assert.Error(t, err)
}

func TestInsert(t *testing.T) {
	service := NewUserService(NewMockUserRepository())

	tests := []struct {
		Name  string
		Error bool
		User  User
	}{
		{
			"Empty",
			true,
			User{},
		},
		{
			"Email Only",
			true,
			User{Email: "test@test.com"},
		},
		{
			"Display Name Only",
			true,
			User{DisplayName: "Test Testerly"},
		},
		{
			"Display Name and Email",
			false,
			User{DisplayName: "Test Testerly", Email: "test@test.com"},
		},
		{
			"Email Takeover",
			true,
			User{DisplayName: "Naughty Tester", Email: "test@test.com"},
		},
		{
			"New User",
			false,
			User{DisplayName: "New User", Email: "new@test.com"},
		},
		{
			"Existing ID",
			true,
			User{Id: primitive.NewObjectID(), DisplayName: "ID", Email: "id@test.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := service.Insert(&test.User)
			if test.Error {
				t.Logf("%s: %v", test.Name, err)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	service := NewUserService(NewMockUserRepository())

	user1 := User{DisplayName: "User One", Email: "one@test.com"}
	assert.NoError(t, service.Insert(&user1))

	user2 := User{DisplayName: "User Two", Email: "two@test.com"}
	assert.NoError(t, service.Insert(&user2))

	t.Run("No ID", func(t *testing.T) {
		user := User{DisplayName: "User NoID", Email: "noid@test.com"}
		assert.Error(t, service.Update(&user))
	})

	t.Run("Simple Update", func(t *testing.T) {
		newName := "One User"
		user1.DisplayName = newName
		assert.NoError(t, service.Update(&user1))
		check, err := service.GetById(user1.Id)
		assert.NoError(t, err)
		assert.Equal(t, newName, check.DisplayName)
	})

	t.Run("Fail Validation", func(t *testing.T) {
		user1.DisplayName = ""
		assert.Error(t, service.Update(&user1))
	})
}
