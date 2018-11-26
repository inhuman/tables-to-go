package database

import (
	"fmt"
	"strings"

	"github.com/inhuman/tables-to-go/src/settings"
	// postgres database driver
	_ "github.com/lib/pq"
)

// postgresql implemenmts the Database interface with help of generalDatabase
type postgresql struct {
	*generalDatabase
}

// DSN creates the DSN String to connect to this database
func (pg *postgresql) DSN(settings *settings.Settings) string {
	return fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		settings.Host, settings.Port, settings.User, settings.DbName, settings.Pswd)
}

// Connect connects to the database by the given data source name (dsn) of the concrete database
func (pg *postgresql) Connect() error {
	return pg.connect(pg.DSN(pg.settings))
}

// GetTables gets all tables for a given schema by name
func (pg *postgresql) GetTables() (tables []*Table, err error) {

	err = pg.db.Select(&tables, `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_type = 'BASE TABLE'
		AND table_schema = $1
		ORDER BY table_name
	`, pg.settings.Schema)

	if pg.settings.Verbose {
		if err != nil {
			fmt.Println("> Error at GetTables()")
			fmt.Printf("> schema: %q\r\n", pg.settings.Schema)
		}
	}

	return tables, err
}

// PrepareGetColumnsOfTableStmt prepares the statement for retrieving the columns of a specific table for a given database
func (pg *postgresql) PrepareGetColumnsOfTableStmt() (err error) {

	pg.getColumnsOfTableStmt, err = pg.db.Preparex(`
		SELECT
			ic.ordinal_position,
			ic.column_name,
			ic.data_type,
			ic.column_default,
			ic.is_nullable,
			ic.character_maximum_length,
			ic.numeric_precision,
			itc.constraint_name,
			itc.constraint_type
		FROM information_schema.columns AS ic
			LEFT JOIN information_schema.key_column_usage AS ikcu ON ic.table_name = ikcu.table_name
			AND ic.table_schema = ikcu.table_schema
			AND ic.column_name = ikcu.column_name
			LEFT JOIN information_schema.table_constraints AS itc ON ic.table_name = itc.table_name
			AND ic.table_schema = itc.table_schema
			AND ikcu.constraint_name = itc.constraint_name
		WHERE ic.table_name = $1
		AND ic.table_schema = $2
		ORDER BY ic.ordinal_position
	`)

	return err
}

// GetColumnsOfTable executes the statement for retrieving the columns of a specific table in a given schema
func (pg *postgresql) GetColumnsOfTable(table *Table) (err error) {

	pg.getColumnsOfTableStmt.Select(&table.Columns, table.Name, pg.settings.Schema)

	if pg.settings.Verbose {
		if err != nil {
			fmt.Printf("> Error at GetColumnsOfTable(%v)\r\n", table.Name)
			fmt.Printf("> schema: %q\r\n", pg.settings.Schema)
		}
	}

	return err
}

// IsPrimaryKey checks if column belongs to primary key
func (pg *postgresql) IsPrimaryKey(column Column) bool {
	return strings.Contains(column.ConstraintType.String, "PRIMARY KEY")
}

// IsAutoIncrement checks if column is a serial column
func (pg *postgresql) IsAutoIncrement(column Column) bool {
	return strings.Contains(column.DefaultValue.String, "nextval")
}

// GetStringDatatypes returns the string datatypes for the postgre database
func (pg *postgresql) GetStringDatatypes() []string {
	return []string{
		"character varying",
		"varchar",
		"character",
		"char",
	}
}

// IsString returns true if colum is of type string for the postgre database
func (pg *postgresql) IsString(column Column) bool {
	return pg.IsStringInSlice(column.DataType, pg.GetStringDatatypes())
}

// GetTextDatatypes returns the text datatypes for the postgre database
func (pg *postgresql) GetTextDatatypes() []string {
	return []string{
		"text",
	}
}

// IsText returns true if colum is of type text for the postgre database
func (pg *postgresql) IsText(column Column) bool {
	return pg.IsStringInSlice(column.DataType, pg.GetTextDatatypes())
}

// GetIntegerDatatypes returns the integer datatypes for the postgre database
func (pg *postgresql) GetIntegerDatatypes() []string {
	return []string{
		"smallint",
		"integer",
		"bigint",
		"smallserial",
		"serial",
		"bigserial",
	}
}

// IsInteger returns true if colum is of type integer for the postgre database
func (pg *postgresql) IsInteger(column Column) bool {
	return pg.IsStringInSlice(column.DataType, pg.GetIntegerDatatypes())
}

// GetFloatDatatypes returns the float datatypes for the postgre database
func (pg *postgresql) GetFloatDatatypes() []string {
	return []string{
		"numeric",
		"decimal",
		"real",
		"double precision",
	}
}

// IsFloat returns true if colum is of type float for the postgre database
func (pg *postgresql) IsFloat(column Column) bool {
	return pg.IsStringInSlice(column.DataType, pg.GetFloatDatatypes())
}

// GetTemporalDatatypes returns the temporal datatypes for the postgre database
func (pg *postgresql) GetTemporalDatatypes() []string {
	return []string{
		"time",
		"timestamp",
		"time with time zone",
		"timestamp with time zone",
		"time without time zone",
		"timestamp without time zone",
		"date",
	}
}

// IsTemporal returns true if colum is of type temporal for the postgre database
func (pg *postgresql) IsTemporal(column Column) bool {
	return pg.IsStringInSlice(column.DataType, pg.GetTemporalDatatypes())
}
