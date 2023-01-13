package gen

import (
	"io/fs"

	"github.com/stephenafamo/bob/gen/drivers"
	"github.com/stephenafamo/bob/orm"
)

// Config for the running of the commands
type Config struct {
	// Additional struct tags to generate (including json, yaml and toml)
	Tags []string `yaml:"tags"`
	// Disable generating factory for models.
	NoFactory bool `yaml:"no_factory"`
	// Disable generated go test files
	NoTests bool `yaml:"no_tests"`
	// Disable back referencing in the loaded relationship structs
	NoBackReferencing bool `yaml:"no_back_referencing"`
	// Delete the output folder (rm -rf) before generation to ensure sanity
	Wipe bool `yaml:"wipe"`
	// Decides the casing for go structure tag names. camel, title or snake (default snake)
	StructTagCasing string `yaml:"struct_tag_casing"`
	// Relationship struct tag name
	RelationTag string `yaml:"relation_tag"`
	// List of column names that should have tags values set to '-' (ignored during parsing)
	TagIgnore []string `yaml:"tag_ignore"`

	Aliases       Aliases       `yaml:"aliases"`
	Replacements  []Replace     `yaml:"replacements"`
	Relationships relationships `yaml:"relationships"`
	Inflections   Inflections   `yaml:"inflections"`

	// Customize the generator name in the top level comment of generated files
	// >>   Code generated by **GENERATOR NAME**. DO NOT EDIT.
	// defaults to "BobGen [plugin] [version]"
	Generator string `yaml:"generator" toml:"generator" json:"generator"`
}

type relationships = map[string][]orm.Relationship

type Output struct {
	PkgName   string  `yaml:"pkg_name" toml:"pkg_name" json:"pkg_name"`
	OutFolder string  `yaml:"out_folder" toml:"out_folder" json:"out_folder"`
	Templates []fs.FS `yaml:"-" toml:"-" json:"-"`

	templates     *templateList
	testTemplates *templateList
}

// Replace replaces a column type with something else
type Replace struct {
	Tables  []string       `yaml:"tables" toml:"tables" json:"tables"`
	Match   drivers.Column `yaml:"match" toml:"match" json:"match"`
	Replace drivers.Column `yaml:"replace" toml:"replace" json:"replace"`
}

type Inflections struct {
	Plural        map[string]string `yaml:"plural"`
	PluralExact   map[string]string `yaml:"plural_exact"`
	Singular      map[string]string `yaml:"singular"`
	SingularExact map[string]string `yaml:"singular_exact"`
	Irregular     map[string]string `yaml:"irregular"`
}
