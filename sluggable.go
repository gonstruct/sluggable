package sluggable

import (
	"fmt"
	"strconv"
	"strings"
)

type Sluggable struct {
	options *options
}

func New(options ...sluggableOption) *Sluggable {
	opts := getDefaultOptions()
	for _, option := range options {
		option(opts)
	}

	return &Sluggable{options: opts}
}

//nolint:cyclop,funlen
func (s *Sluggable) Generate(db contextExecutor, value string, options ...sluggableOption) (string, error) {
	opts := s.options
	for _, option := range options {
		option(opts)
	}

	if len(opts.tableName) == 0 {
		return "", fmt.Errorf("[sluggable] table name cannot be empty")
	}

	slug := opts.method(value, opts.seperator)

	sql := `SELECT "id", "{column}" FROM "{table}" WHERE ("{column}" = $1 OR "{column}" LIKE $2)`

	params := []any{slug, fmt.Sprint(slug, opts.seperator, "%")}
	for whereSql, args := range opts.wheres {
		normalizedSql := whereSql

		for i := 0; i < len(args); i++ {
			placeholder := fmt.Sprintf("$%d", len(params)+1)
			normalizedSql = strings.ReplaceAll(normalizedSql, "?", placeholder)
			params = append(params, args[i])
		}

		sql += fmt.Sprintf(" AND (%s)", normalizedSql)
		params = append(params, args...)
	}

	sql = strings.ReplaceAll(sql, "{table}", opts.tableName)
	sql = strings.ReplaceAll(sql, "{column}", opts.columnName)

	rows, err := db.Query(sql, params...)
	if err != nil {
		return "", fmt.Errorf("[sluggable] failed to query sluggable: %w", err)
	}
	defer rows.Close()

	simularList := make(map[string]string)
	for rows.Next() {
		var idValue string
		var slugValue string
		if err := rows.Scan(&idValue, &slugValue); err != nil {
			return "", fmt.Errorf("[sluggable] failed to scan sluggable value: %w", err)
		}
		simularList[idValue] = slugValue
	}

	if len(simularList) == 0 {
		return slug, nil
	}

	if opts.identifier != "" {
		if existingSlug, exists := simularList[opts.identifier]; exists {
			if existingSlug == slug || existingSlug == "" || strings.HasPrefix(existingSlug, slug) {
				return existingSlug, nil
			}
		}
	}

	latestSuffix := 0
	for _, simular := range simularList {
		suffix := strings.TrimPrefix(simular, fmt.Sprint(slug, opts.seperator))
		suffixAsNumber, err := strconv.Atoi(suffix)
		if err != nil {
			continue
		}

		if suffixAsNumber > latestSuffix {
			latestSuffix = suffixAsNumber
		}
	}

	if latestSuffix > 0 {
		return fmt.Sprint(slug, opts.seperator, latestSuffix+1), nil
	}

	return fmt.Sprint(slug, opts.seperator, opts.firstUniqueSuffix), nil
}

func Generate(db contextExecutor, value string, options ...sluggableOption) (string, error) {
	if _global == nil {
		_global = New()
	}

	return _global.Generate(db, value, options...)
}
