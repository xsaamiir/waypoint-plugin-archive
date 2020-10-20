package builder

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/mholt/archiver/v3"
)

type BuilderConfig struct {
	// Sources is a list of files and/or directories to package inside the archive.
	Sources []string `hcl:"sources"`
	// OutputName is the name of the archive file.
	OutputName string `hcl:"output_name"`
	// OverwriteExisting indicates whether to overwrite the existing file;
	// if false, an error is returned if the file exists.
	OverwriteExisting bool `hcl:"overwrite_existing,optional"`
}

type Builder struct {
	config BuilderConfig
}

// Documentation implements Documented.
func (b *Builder) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&BuilderConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Archive")

	doc.Example(`
build {
  use "archive" {
      sources = ["./src", "./public", "./package.json"]      
      output_name = "webapp.zip"
      overwrite_existing = true
  }
}
`)

	doc.Input("component.Source")
	doc.Output("archive.Archive")

	_ = doc.SetField(
		"sources",
		"a list of files and/or directories to package inside the archive",
		docs.Summary(
			"The list of files and/or directoires to package inside the archive. "+
				"The sources should be relative to the path of the application being built. ",
			"Ex: `/path/to/project/app/`",
		),
	)

	_ = doc.SetField("output_name", "the name of the archive file", docs.Summary())

	_ = doc.SetField(
		"overwrite_existing",
		"if false, an error is returned if the file exists.",
		docs.Summary("Whether to overwrite the existing file; if false, an error is returned if the file exists."),
	)

	return doc, nil
}

// Config implements Configurable.
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// ConfigSet implements ConfigurableNotify.
func (b *Builder) ConfigSet(config interface{}) error {
	c, ok := config.(*BuilderConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit.
		return fmt.Errorf("Expected *BuildConfig as parameter")
	}

	// validate the config
	sources := c.Sources
	if len(sources) == 0 {
		return errors.New("Sources can't be empty, please provide the path to at least one file or directory")
	}

	return nil
}

// BuildFunc implements Builder.
func (b *Builder) BuildFunc() interface{} {
	// return a function which will be called by Waypoint
	return b.build
}

// A BuildFunc does not have a strict signature, you can define the parameters
// you need based on the Available parameters that the Waypoint SDK provides.
// Waypoint will automatically inject parameters as specified
// in the signature at run time.
//
// Available input parameters:
// - context.Context
// - *component.Source
// - *component.JobInfo
// - *component.DeploymentConfig
// - *datadir.Project
// - *datadir.App
// - *datadir.Component
// - hclog.Logger
// - terminal.UI
// - *component.LabelSet
//
// The output parameters for BuildFunc must be a Struct which can
// be serialzied to Protocol Buffers binary format and an error.
// This Output Value will be made available for other functions
// as an input parameter.
// If an error is returned, Waypoint stops the execution flow and
// returns an error to the user.
func (b *Builder) build(source *component.Source, logger hclog.Logger, ui terminal.UI) (*Archive, error) {
	logger.Trace("creating a new archive", "config", b.config)

	u := ui.Status()
	defer u.Close()
	u.Update("Creating archive")

	cwd, err := os.Getwd()
	if err != nil {
		u.Step(terminal.StatusError, "Error getting current working directory")
		return nil, err
	}

	sourcePath := source.Path
	sources := b.config.Sources

	for i, src := range sources {
		sources[i] = path.Join(sourcePath, src)
	}

	outputName := b.config.OutputName
	outputPath := path.Join(cwd, sourcePath, outputName)

	zip := archiver.NewZip()
	zip.OverwriteExisting = b.config.OverwriteExisting

	err = zip.Archive(sources, outputPath)
	if err != nil {
		u.Step(terminal.StatusError, "Archive failed")
		return nil, err
	}

	u.Step(terminal.StatusOK, "Archive saved to '"+outputPath+"'")

	return &Archive{OutputPath: outputPath}, nil
}
