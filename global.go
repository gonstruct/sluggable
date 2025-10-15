package sluggable

import (
	slugify "github.com/gosimple/slug"
)

var _global *Sluggable

const (
	excludeDeletedWhere = `"deleted_at" IS NULL`
)

func getDefaultOptions() *options {
	return &options{
		method: func(value, seperator string) string {
			return slugify.MakeLang(value, "en")
		},
		seperator:         "-",
		tableName:         "",
		columnName:        "slug",
		firstUniqueSuffix: 2,
		wheres: map[string][]any{
			excludeDeletedWhere: {},
		},
	}
}

func Configure(options ...sluggableOption) {
	if _global == nil {
		_global = New(options...)
		return
	}

	for _, option := range options {
		option(_global.options)
	}
}
