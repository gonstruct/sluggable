package main

import (
	"database/sql"

	"github.com/gonstruct/sluggable"
)

type DatabaseModel struct {
	ID           string       `db:"id"`
	Name         string       `db:"name"`
	Slug         string       `db:"slug"`
	CustomColumn int          `db:"custom_column"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

//nolint:gochecknoinits
func init() {
	sluggable.Configure(
		sluggable.WithDeleted(), // Include deleted records
		sluggable.WithWhere(`custom_column = ? AND other_column = ? `, 1, 2),
	)
}

func (d *DatabaseModel) OnCreate(db *sql.DB) (err error) {
	d.Slug, err = sluggable.Generate(db, d.Name, sluggable.WithTableName("database_models"))

	return err
}

func (d *DatabaseModel) OnUpdate(db *sql.DB) (err error) {
	// Or use a custom client
	myCustomSluggable := sluggable.New( // Add this in your service
		sluggable.WithDeleted(), // Include deleted records
		sluggable.WithWhere(`custom_column = ? AND other_column = ? `, 1, 2),
	)

	d.Slug, err = myCustomSluggable.Generate(db, d.Name,
		sluggable.WithTableName("database_models"),
		sluggable.WithIdentifier(d.ID), // Exclude current record when checking for existing slugs
	)

	return err
}
