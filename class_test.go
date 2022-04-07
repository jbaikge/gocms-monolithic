package gocms

import (
	"fmt"
	"testing"
	"time"

	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ ClassRepository = mockClassRepository{}

type mockClassRepository struct {
	byId   map[primitive.ObjectID]Class
	bySlug map[string]Class
}

func NewMockClassRepository() mockClassRepository {
	return mockClassRepository{
		byId:   make(map[primitive.ObjectID]Class),
		bySlug: make(map[string]Class),
	}
}

func (r mockClassRepository) DeleteClass(id primitive.ObjectID) (err error) {
	class, ok := r.byId[id]
	if !ok {
		// Silent failure
		return
	}
	delete(r.byId, id)
	delete(r.bySlug, class.Slug)
	return
}

func (r mockClassRepository) GetAllClasses() (all []Class, err error) {
	all = make([]Class, 0, len(r.byId))
	for _, class := range r.byId {
		all = append(all, class)
	}
	return
}

func (r mockClassRepository) GetClassById(id primitive.ObjectID) (class Class, err error) {
	class, ok := r.byId[id]
	if !ok {
		err = fmt.Errorf("class not found: %s", id)
	}
	return
}

func (r mockClassRepository) GetClassBySlug(slug string) (class Class, err error) {
	class, ok := r.bySlug[slug]
	if !ok {
		err = fmt.Errorf("class not found: %s", slug)
	}
	return
}

func (r mockClassRepository) InsertClass(class *Class) (err error) {
	class.Id = primitive.NewObjectID()
	r.byId[class.Id] = *class
	r.bySlug[class.Slug] = *class
	return
}

func (r mockClassRepository) UpdateClass(class *Class) (err error) {
	class.Updated = time.Now()
	if err = r.DeleteClass(class.Id); err != nil {
		return
	}
	r.byId[class.Id] = *class
	r.bySlug[class.Slug] = *class
	return
}

func TestClassService(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		service := NewClassService(NewMockClassRepository())

		classes := []*Class{
			{Name: "Test", Slug: "test1"},
			{Name: "Test", Slug: "test2"},
		}
		for _, c := range classes {
			assert.NoError(t, service.Insert(c))
		}

		all, err := service.All()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(all))
	})

	t.Run("GetById", func(t *testing.T) {
		service := NewClassService(NewMockClassRepository())

		class := Class{Name: "Test", Slug: "test"}
		assert.NoError(t, service.Insert(&class))

		check, err := service.GetById(class.Id)
		assert.NoError(t, err)
		assert.Equal(t, class.Id, check.Id)
	})

	t.Run("GetBySlug", func(t *testing.T) {
		service := NewClassService(NewMockClassRepository())

		class := Class{Name: "Test", Slug: "test"}
		assert.NoError(t, service.Insert(&class))

		check, err := service.GetBySlug(class.Slug)
		assert.NoError(t, err)
		assert.Equal(t, class.Id, check.Id)
	})

	t.Run("Insert", func(t *testing.T) {
		service := NewClassService(NewMockClassRepository())

		tests := []struct {
			Name  string
			Error bool
			Class Class
		}{
			{
				"Slug Only",
				true,
				Class{Slug: "test"},
			},
			{
				"Name Only",
				true,
				Class{Name: "Test"},
			},
			{
				"Slug & Name",
				false,
				Class{Name: "Test", Slug: "test"},
			},
			{
				"Same Slug",
				true,
				Class{Name: "Test", Slug: "test"},
			},
			{
				"Pre-ID",
				true,
				Class{Id: primitive.NewObjectID(), Name: "Test", Slug: "pre_id"},
			},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				err := service.Insert(&test.Class)
				if test.Error {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		service := NewClassService(NewMockClassRepository())

		t.Run("No ID", func(t *testing.T) {
			class := Class{Name: "No ID", Slug: "no_id"}
			assert.Error(t, service.Update(&class))
		})

		banana := Class{Name: "Banana", Slug: "banana"}
		assert.NoError(t, service.Insert(&banana))

		orange := Class{Name: "Orange", Slug: "orange"}
		assert.NoError(t, service.Insert(&orange))

		t.Run("Name Update", func(t *testing.T) {
			orange.Name = "Tangerine"
			assert.NoError(t, service.Update(&orange))

			orangeTest, err := service.GetById(orange.Id)
			assert.NoError(t, err)
			assert.Equal(t, orange.Name, orangeTest.Name)
		})

		t.Run("Blank Slug", func(t *testing.T) {
			orange.Slug = ""
			assert.Error(t, service.Update(&orange))
		})

		t.Run("Slug Takeover", func(t *testing.T) {
			banana.Slug = "orange"
			assert.Error(t, service.Update(&banana))
		})
	})

	t.Run("Delete", func(t *testing.T) {
		service := NewClassService(NewMockClassRepository())

		class := Class{Name: "Test", Slug: "test"}
		assert.NoError(t, service.Insert(&class))
		assert.NoError(t, service.Delete(class))
		// Do it once more to make sure it fails silently
		assert.NoError(t, service.Delete(class))

		_, err := service.GetById(class.Id)
		assert.Error(t, err)
	})
}
