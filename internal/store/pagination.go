package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedQuery struct {
	Limit  int        `json:"limit" validate:"gte=1,lte=20"`
	Offset int        `json:"offset" validate:"gte=0"`
	Sort   string     `json:"sort" validate:"oneof=ASC DESC"`
	Search string     `json:"search" validate:"max=100"`
	Since  *time.Time `json:"since"`
	Until  *time.Time `json:"until"`
	Tags   []string   `json:"tags" validate:"max=5"`
}

func (query *PaginatedQuery) Parse(r *http.Request) error {

	queryValues := r.URL.Query()

	limit := queryValues.Get("limit")
	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		query.Limit = parsedLimit
	}

	offset := queryValues.Get("offset")
	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil {
			return err
		}
		query.Offset = parsedOffset
	}

	query.Search = queryValues.Get("search")
	query.Until = parseTime(queryValues.Get("until"))
	query.Since = parseTime(queryValues.Get("since"))

	tags := queryValues.Get("tags")
	if tags != "" {
		query.Tags = strings.Split(tags, ",")
	}

	return nil
}

func parseTime(str string) *time.Time {
	parsedTime, err := time.Parse("2006-01-02", str)
	if err != nil {
		return nil
	}
	return &parsedTime
}
