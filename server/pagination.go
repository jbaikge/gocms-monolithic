package server

import (
	"fmt"
	"math"
)

type PaginationLink struct {
	Label    string
	Page     int64
	Disabled bool
	Active   bool
}

type Pagination struct {
	Page          int64
	PerPage       int64
	Total         int64
	Shoulder      int64
	PreviousLabel string
	NextLabel     string
}

func NewPagination(page, perPage, total int64) Pagination {
	return Pagination{
		Page:          page,
		PerPage:       perPage,
		Total:         total,
		Shoulder:      2,
		PreviousLabel: "Previous",
		NextLabel:     "Next",
	}
}

func (p Pagination) Links() (links []PaginationLink) {
	pages := p.windowPages()
	if len(pages) == 0 {
		return
	}

	lastPage := p.pageCount()
	links = make([]PaginationLink, 0, len(pages)+4)

	// Add previous
	links = append(links, PaginationLink{
		Label:    p.PreviousLabel,
		Page:     int64(math.Max(1, float64(p.Page-1))),
		Disabled: p.Page == 1,
	})

	// Add page 1
	if pages[0] > 1 {
		links = append(links, PaginationLink{
			Label:  "1",
			Page:   1,
			Active: p.Page == 1,
		})
	}

	for _, page := range pages {
		links = append(links, PaginationLink{
			Label:  fmt.Sprint(page),
			Page:   page,
			Active: p.Page == page,
		})
	}

	// Add last page
	if pages[len(pages)-1] < lastPage {
		links = append(links, PaginationLink{
			Label:  fmt.Sprint(lastPage),
			Page:   lastPage,
			Active: p.Page == lastPage,
		})
	}

	// Add next
	links = append(links, PaginationLink{
		Label:    p.NextLabel,
		Page:     int64(math.Min(float64(lastPage), float64(p.Page+1))),
		Disabled: p.Page == lastPage,
	})

	return
}

// Window pages are the pages between 1 and the last page number, centered on
// the current page.
//
// This is how you get pagination that looks like:
// 1 ... 3 4 5 6 7 ... 10
//
// The shoulder defines how many pages on each side of the current page to
// display.
func (p Pagination) windowPages() (pages []int64) {
	lastPage := p.pageCount()
	if lastPage <= 1 {
		return
	}

	// The easy one: all pages fit inside the window
	// Multiply by 2 to handle each shoulder
	// Add three to handle page 1, current page, and last page
	if p.Shoulder*2+3 >= lastPage {
		pages = make([]int64, lastPage)
		for i := range pages {
			pages[i] = int64(i + 1)
		}
		return
	}

	// The tough one: create a window centered around the current page
	min := float64(2)
	max := float64(lastPage - 1)
	page := float64(p.Page)
	shoulder := float64(p.Shoulder)
	rangeMin := math.Min(max-shoulder*2, math.Max(min, page-shoulder))
	rangeMax := math.Max(min+shoulder*2, math.Min(max, page+shoulder))
	pages = make([]int64, int(rangeMax-rangeMin))
	for i := range pages {
		pages[i] = int64(i) + int64(rangeMin)
	}

	return
}

func (p Pagination) pageCount() int64 {
	perPage := float64(p.PerPage)
	total := float64(p.Total)
	return int64(math.Ceil(total / perPage))
}
