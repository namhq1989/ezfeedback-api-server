package main

import (
	"os"
	"slices"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
)

var excludeTables = []string{"schema_version", "spatial_ref_sys"}

func main() {
	conn, err := parsePostgresConnectionString(os.Getenv("POSTGRES_CONN"))
	if err != nil {
		panic(err)
	}

	if err = postgres.Generate(
		"internal/database/gen",
		*conn,
		template.Default(postgres2.Dialect).
			UseSchema(func(schema metadata.Schema) template.Schema {
				return template.DefaultSchema(schema).
					UseSQLBuilder(template.DefaultSQLBuilder().
						UseView(func(table metadata.Table) template.TableSQLBuilder {
							return template.TableSQLBuilder{
								Skip: true,
							}
						}).
						UseTable(func(table metadata.Table) template.TableSQLBuilder {
							if slices.Contains(excludeTables, table.Name) {
								return template.TableSQLBuilder{Skip: true}
							}

							return template.DefaultTableSQLBuilder(table)
						}),
					).
					UseModel(template.DefaultModel().
						UseView(func(table metadata.Table) template.ViewModel {
							return template.ViewModel{
								Skip: true,
							}
						}).
						UseTable(func(table metadata.Table) template.TableModel {
							if slices.Contains(excludeTables, table.Name) {
								return template.TableModel{Skip: true}
							}
							return template.DefaultTableModel(table).
								UseField(func(column metadata.Column) template.TableModelField {
									field := template.DefaultTableModelField(column)

									if schema.Name == "public" {
										switch table.Name {
										case "words":
											field = word(field, column)
										case "word_news":
											field = wordNews(field, column)
										}
									}

									return field
								})
						}),
					)
			}),
	); err != nil {
		panic(err)
	}
}

func word(field template.TableModelField, column metadata.Column) template.TableModelField {
	switch column.Name {
	case "parts_of_speech":
		field.Type = template.NewType(database.ArrayString{})
	}

	return field
}

func wordNews(field template.TableModelField, column metadata.Column) template.TableModelField {
	switch column.Name {
	case "categories":
		field.Type = template.NewType(database.ArrayString{})
	}

	return field
}
