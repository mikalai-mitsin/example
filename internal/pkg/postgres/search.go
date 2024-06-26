package postgres

import (
	"fmt"
	"strings"
)

type Search struct {
	Lang   string
	Fields []string
	Query  string
}

// nolint:stylecheck
func (s Search) ToSql() (string, []interface{}, error) {
	if s.Lang == "" {
		s.Lang = "english"
	}
	vector := "to_tsvector('%s', %s)"
	vector = fmt.Sprintf(vector, s.Lang, strings.Join(s.Fields, " || ' ' || "))
	query := "plainto_tsquery('%s', ?)"
	query = fmt.Sprintf(query, s.Lang)
	args := []interface{}{s.Query}
	return fmt.Sprintf("%s @@ %s", vector, query), args, nil
}
