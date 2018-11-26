package tagger

import (
	"fmt"
	"github.com/inhuman/tables-to-go/src/database"
	"github.com/inhuman/tables-to-go/src/settings"
	"github.com/stoewer/go-strcase"
	"strings"
)

// Tagger interface for types of struct-tages
type Tagger interface {
	GenerateTag(db database.Database, column database.Column, settings *settings.Settings) string
}

// DbTag is the standard "db"-tag
type DbTag string

// GenerateTag for DbTag to satisfy the Tagger interface
func (t *DbTag) GenerateTag(db database.Database, column database.Column, settings *settings.Settings) string {

	tag := `db:"` + column.Name + `"`

	if settings.OutputFormatTag == "o" {
		tag = `db:"` + strcase.SnakeCase(column.Name) + `"`
	}

	return tag
}

// StblTag represents the Masterminds/structable "stbl"-tag
type StblTag string

// GenerateTag for StblTag to satisfy the Tagger interface
func (t *StblTag) GenerateTag(db database.Database, column database.Column, settings *settings.Settings) string {

	isPk := ""
	if db.IsPrimaryKey(column) {
		isPk = ",PRIMARY_KEY"
	}

	isAutoIncrement := ""
	if db.IsAutoIncrement(column) {
		isAutoIncrement = ",SERIAL,AUTO_INCREMENT"
	}

	tag := `stbl:"` + column.Name + isPk + isAutoIncrement + `"`

	if settings.OutputFormatTag == "o" {
		tag = `stbl:"` + strcase.SnakeCase(column.Name) + isPk + isAutoIncrement + `"`
	}

	return tag

}

// SQLTag is the experimental "sql"-tag
type SQLTag string

// GenerateTag for SQLTag to satisfy the Tagger interface
func (t *SQLTag) GenerateTag(db database.Database, column database.Column, settings *settings.Settings) string {

	colType := ""
	characterMaximumLength := ""
	if db.IsString(column) && column.CharacterMaximumLength.Valid {
		characterMaximumLength = fmt.Sprintf("(%v)", column.CharacterMaximumLength.Int64)
	}

	colType = fmt.Sprintf("type:%v%v;", column.DataType, characterMaximumLength)

	isNullable := ""
	if !db.IsNullable(column) {
		isNullable = "not null;"
	}

	// TODO size:###
	// TODO unique, key, index, ...

	tag := colType + isNullable
	tag = strings.TrimSuffix(tag, ";")

	return `sql:"` + tag + `"`
}
