package store

import (
	"net/http"
	"strconv"
	"strings"
)


type PaginatedFeed struct {
	Limit 	int `json:"limit" validate:"gte=1,lte=20"`
	Offset	int `json:"offset" validate:"gte=0"`
	Sort 	string `json:"sort" validate:"oneof=asc desc"`
	Search 	string `json:"search" validate:"max=100"`
	Tags	[]string `json:"tags" validate:"max=5"`
}

func (p PaginatedFeed) Parse(r *http.Request) (PaginatedFeed, error) {
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return p, err
		}
		p.Limit = limitInt
	}

	offset := r.URL.Query().Get("offset")
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return p, err
		}
		p.Offset = offsetInt
	}

	sort := r.URL.Query().Get("sort")
	if sort != "" {
		p.Sort = sort
	}

	search := r.URL.Query().Get("search")
	if search != "" {
		p.Search = search
	}

	tags := r.URL.Query().Get("tags")
	// split tags by comma
	if tags != "" {
		p.Tags = strings.Split(tags, ",")
	}


	return p, nil
}
