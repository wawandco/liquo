package liquo

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/go-xmlfmt/xmlfmt"
	"github.com/gobuffalo/flect"
	"github.com/spf13/pflag"
	"github.com/wawandco/liquo/internal/log"
	"github.com/wawandco/ox/plugins/core"
)

var (
	ErrNameArgMissing = errors.New("name arg missing")
	ErrInvalidName    = errors.New("invalid migration name")
	ErrInvalidPath    = errors.New("invalid path")

	// MigrationTemplate for the migration generator.
	//go:embed templates/migration.xml.tmpl
	migrationTemplate string
)

var (
	// Ensuring we're building a plugin
	_ core.Plugin = (*Generator)(nil)
	// Ensuring the plugin is a flagparser
	_ core.FlagParser = (*Generator)(nil)
)

// Generator for liquibase SQL migrations, it generates xml liquibase
// for SQL in the root + basedir folder. It uses the argument passed
// to determine both the name of the migration and the destination.
// Some examples are:
// - "ox generate migration name" generates [timestamp]-name.xml
// - "ox generate migration folder/name" generates folder/[timestamp]-name.xml
// - "ox generate migration name --base migrations" generates migrations/[timestamp]-name.xml
type Generator struct {
	// mockTimestamp is used for testing purposes, it would replace the
	// timestamp at the beggining of the migration name.
	mockTimestamp string

	// Basefolder for the migrations, if a path is passed, then we will append that
	// path to the baseFolder when generating the migration.
	baseFolder string

	flags *pflag.FlagSet
}

// Name is the name used to identify the generator and also
// the plugin
func (g Generator) Name() string {
	return "liquo/generate-migration"
}

// Name is the name used to identify the generator and also
// the plugin
func (g Generator) InvocationName() string {
	return "migration"
}

// Generate a new migration based on the passed args. This needs at least 3
// args since the 3rd arg will be used by the generator to build the name of
// the migration.
func (g Generator) Generate(ctx context.Context, root string, args []string) error {
	if len(args) < 3 {
		return ErrNameArgMissing
	}

	path, err := g.generateFile(args)
	if err != nil {
		return err
	}

	log.Infof("migration generated in %v", path)
	err = g.addToChangelog(root, path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil {
		log.Infof("migration added to the changelog")
	}

	return nil
}

func (g Generator) addToChangelog(root, path string) error {
	changelog := filepath.Join(root, "migrations", "changelog.xml")
	original, err := ioutil.ReadFile(changelog)
	if err != nil {
		return err
	}

	statement := fmt.Sprintf(`<include file="%s" />`, path)
	result := strings.Replace(string(original), `</databaseChangeLog>`, statement+"</databaseChangeLog>", 1)
	result = xmlfmt.FormatXML(result, "", "\t")
	parts := strings.Split(result, "\n")
	result = strings.Join(parts[1:], "\n")
	err = ioutil.WriteFile(changelog, []byte(result), 0777)
	if err != nil {
		return err
	}

	return nil
}

func (g Generator) generateFile(args []string) (string, error) {
	timestamp := time.Now().UTC().Format("20060102150405")
	if g.mockTimestamp != "" {
		timestamp = g.mockTimestamp
	}

	filename, err := g.composeFilename(args[2], timestamp)
	if err != nil {
		return "", err
	}

	path := g.baseFolder
	if dir := filepath.Dir(args[2]); dir != "." {
		path = filepath.Join(g.baseFolder, dir)
	}

	path = filepath.Join(path, filename)
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return path, err
	}

	tmpl, err := template.New("migration-template").Parse(migrationTemplate)
	if err != nil {
		return path, err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, strings.ReplaceAll(filename, ".xml", ""))
	if err != nil {
		return path, err
	}

	return path, ioutil.WriteFile(path, tpl.Bytes(), 0655)
}

// composeFilename from the passed arg and timestamp, if the passed path is
// a dot (.) or a folder "/" then it will return ErrInvalidName.
func (g Generator) composeFilename(passed, timestamp string) (string, error) {
	name := filepath.Base(passed)
	//Should we check the name here ?
	if name == "." || name == "/" {
		return "", ErrInvalidName
	}

	underscoreName := flect.Underscore(name)
	result := timestamp + "-" + underscoreName + ".xml"

	return result, nil
}

// Parseflags will parse the baseFolder from the --base or -b flag
func (g *Generator) ParseFlags(args []string) {
	g.flags = pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)
	g.flags.StringVarP(&g.baseFolder, "base", "b", "migrations", "destination folder of the generated migration")
	g.flags.Parse(args) //nolint:errcheck,we don't care hence the flag
}

// Flags parsed by the plugin
func (g *Generator) Flags() *pflag.FlagSet {
	return g.flags
}
