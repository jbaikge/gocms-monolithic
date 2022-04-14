package server

import (
	"testing"

	"github.com/zeebo/assert"
)

func TestPagination(t *testing.T) {
	tests := []struct {
		Name       string
		Pagination Pagination
		Expect     []PaginationLink
	}{
		{
			"One Page",
			NewPagination(1, 10, 10),
			nil,
		},
		{
			"Two Pages",
			NewPagination(1, 10, 20),
			[]PaginationLink{
				{Label: "Previous", Page: 1, Disabled: true},
				{Label: "1", Page: 1, Active: true},
				{Label: "2", Page: 2},
				{Label: "Next", Page: 2},
			},
		},
		{
			"With Shoulders",
			NewPagination(5, 10, 90),
			[]PaginationLink{
				{Label: "Previous", Page: 4},
				{Label: "1", Page: 1},
				{Label: "3", Page: 3},
				{Label: "4", Page: 4},
				{Label: "5", Page: 5, Active: true},
				{Label: "6", Page: 6},
				{Label: "7", Page: 7},
				{Label: "9", Page: 9},
				{Label: "Next", Page: 6},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.DeepEqual(t, test.Expect, test.Pagination.Links())
		})
	}
}
