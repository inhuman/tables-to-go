package settings

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	// supportedDbTypes represents the supported databases
	supportedDbTypes = map[string]bool{
		"pg":    true,
		"mysql": true,
	}

	// supportedOutputFormats represents the supported output formats
	supportedOutputFormats = map[string]bool{
		"c": true,
		"o": true,
	}

	// dbDefaultPorts maps the database type to the default ports
	dbDefaultPorts = map[string]string{
		"pg":    "5432",
		"mysql": "3306",
	}
)

var GlobalSettings *Settings

// Settings stores the supported settings / command line arguments
type Settings struct {
	Verbose         bool
	DbType          string
	User            string
	Pswd            string
	DbName          string
	Schema          string
	Host            string
	Port            string
	OutputFilePath  string
	OutputFormat    string
	OutputFormatTag string
	PackageName     string
	Prefix          string
	Suffix          string

	TagsNoDb bool

	TagsMastermindStructable       bool
	TagsMastermindStructableOnly   bool
	IsMastermindStructableRecorder bool

	// TODO not implemented yet
	TagsGorm bool

	// experimental
	TagsSQL     bool
	TagsSQLOnly bool
}

// NewSettings constructs settings with default values
func NewSettings() *Settings {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = "."
	}

	return &Settings{
		Verbose:         false,
		DbType:          "pg",
		User:            "postgres",
		Pswd:            "",
		DbName:          "postgres",
		Schema:          "public",
		Host:            "127.0.0.1",
		Port:            "", // left blank -> is automatically determined if not set
		OutputFilePath:  dir,
		OutputFormat:    "c",
		OutputFormatTag: "c",
		PackageName:     "dto",
		Prefix:          "",
		Suffix:          "",

		TagsNoDb: false,

		TagsMastermindStructable:       false,
		TagsMastermindStructableOnly:   false,
		IsMastermindStructableRecorder: false,

		TagsGorm: false,

		TagsSQL:     false,
		TagsSQLOnly: false,
	}
}

// Verify verifies the settings and checks the given output paths
func (settings *Settings) Verify() (err error) {

	if !supportedDbTypes[settings.DbType] {
		return fmt.Errorf("type of database %q not supported! %v", settings.DbType, settings.SupportedDbTypes())
	}

	if !supportedOutputFormats[settings.OutputFormat] {
		return fmt.Errorf("output format %q not supported", settings.OutputFormat)
	}

	if err = settings.verifyOutputPath(); err != nil {
		return err
	}

	if settings.OutputFilePath, err = settings.prepareOutputPath(); err != nil {
		return err
	}

	if settings.Port == "" {
		settings.Port = dbDefaultPorts[settings.DbType]
	}

	if settings.PackageName == "" {
		return fmt.Errorf("name of package can not be empty")
	}

	return err
}

func (settings *Settings) verifyOutputPath() (err error) {

	info, err := os.Stat(settings.OutputFilePath)

	if os.IsNotExist(err) {
		return fmt.Errorf("output file path %q does not exists", settings.OutputFilePath)
	}

	if !info.Mode().IsDir() {
		return fmt.Errorf("output file path %q is not a directory", settings.OutputFilePath)
	}

	return err
}

func (settings *Settings) prepareOutputPath() (outputFilePath string, err error) {
	outputFilePath, err = filepath.Abs(settings.OutputFilePath)
	outputFilePath += string(filepath.Separator)
	return outputFilePath, err
}

// SupportedDbTypes returns a slice of strings as names of the supported database types
func (settings *Settings) SupportedDbTypes() string {
	names := make([]string, len(supportedDbTypes))
	i := 0
	for name := range supportedDbTypes {
		names[i] = name
		i++
	}
	return fmt.Sprintf("%v", names)
}
