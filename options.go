package sluggable

type options struct {
	method    func(value, seperator string) string // Defaults to "slugify"
	seperator string                               // Defaults to "-"

	tableName  string // Empty by default, must be set
	columnName string // Defaults to "slug"

	identifier string // Optional, used to check for existing slugs

	firstUniqueSuffix int // Defaults to 2

	wheres map[string][]any // Optional, used to add additional where clauses
}

type sluggableOption func(*options)

func WithMethod(method func(value, seperator string) string) sluggableOption {
	return func(opts *options) {
		opts.method = method
	}
}

func WithSeperator(seperator string) sluggableOption {
	return func(opts *options) {
		opts.seperator = seperator
	}
}

func WithTableName(tableName string) sluggableOption {
	return func(opts *options) {
		opts.tableName = tableName
	}
}

func WithColumnName(columnName string) sluggableOption {
	return func(opts *options) {
		opts.columnName = columnName
	}
}

func WithFirstUniqueSuffix(suffix int) sluggableOption {
	return func(opts *options) {
		opts.firstUniqueSuffix = suffix
	}
}

func WithIdentifier(identifier string) sluggableOption {
	return func(opts *options) {
		opts.identifier = identifier
	}
}

func WithDeleted() sluggableOption {
	return func(opts *options) {
		delete(opts.wheres, excludeDeletedWhere)
	}
}

func WithWhere(sql string, params ...any) sluggableOption {
	return func(opts *options) {
		opts.wheres[sql] = params
	}
}
