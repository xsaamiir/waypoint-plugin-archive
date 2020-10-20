package builder

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type BuilderConfig struct {
	// Sources is a list of files and/or directories to package inside the archive.
	Sources []string `hcl:"sources"`
	// OutputName is the name of the archive file.
	OutputName string `hcl:"output_name"`
	// OverwriteExisting indicates whether to overwrite the existing file;
	// if false, an error is returned if the file exists.
	OverwriteExisting bool `hcl:"overwrite_existing,optional"`
	// Ignore is a list of files or directories to ignore when reading sources.
	Ignore []string `hcl:"ignore,optional"`
	// CollapseTopLevelFolder indicates whether to include the source directory in
	// the archive or only add its content.
	CollapseTopLevelFolder bool `hcl:"collapse_top_level_folder,optional"`
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
      sources = ["."]      
      output_name = "webapp.zip"
      overwrite_existing = true
      ignore = ["node_modules", "src/docs/README.md"]
	  collapse_top_level_folder = true
  }
}
`)

	doc.Input("component.Source")
	doc.Output("archive.Archive")

	_ = doc.SetField(
		"sources",
		"The list of files and/or directories to package inside the archive",
		docs.Summary(
			"The list of files and/or directoires to package inside the archive. "+
				"The sources should be relative to the path of the application being built. ",
			"Ex: `/path/to/project/app/`",
		),
	)

	_ = doc.SetField("output_name", "The name of the archive file")

	_ = doc.SetField(
		"overwrite_existing",
		"Whether to overwrite the existing file; if false, an error is returned if the file exists.",
		docs.Default("false"),
	)

	_ = doc.SetField(
		"ignore",
		"A list of paths to files and/or directories to ignore while creating the archive. "+
			"The paths should be relative to the folder of the app.",
	)

	_ = doc.SetField(
		"collapse_top_level_folder",
		"Whether to add to the source directory in the archive or only its content; if false, "+
			"the source directory will be included in the archive, if true only the content will be included.",
		docs.Default("false"),
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
// be serialized to Protocol Buffers binary format and an error.
// This Output Value will be made available for other functions
// as an input parameter.
// If an error is returned, Waypoint stops the execution flow and
// returns an error to the user.
func (b *Builder) build(source *component.Source, logger hclog.Logger, ui terminal.UI) (*Archive, error) {
	logger.Debug("creating a new archive", "config", b.config)

	st := ui.Status()
	defer st.Close()

	st.Update("Creating archive")

	cwd, err := os.Getwd()
	if err != nil {
		st.Step(terminal.StatusError, "Error getting current working directory")
		return nil, err
	}

	sourcePath := path.Join(cwd, source.Path)
	sources := b.config.Sources
	ignore := b.config.Ignore
	xsources := make([]string, 0)
	outputPath := path.Join(sourcePath, b.config.OutputName)

	if !b.config.OverwriteExisting && fileExists(outputPath) {
		return nil, errors.New("Output file already exists")
	}

	for _, src := range sources {
		xsrc, err := expandSource(path.Join(sourcePath, src), ignore)
		if err != nil {
			st.Step(terminal.StatusError, "Error expanding source")
			return nil, err
		}

		xsources = append(xsources, xsrc...)
	}

	basePath := cwd
	if b.config.CollapseTopLevelFolder {
		basePath = sourcePath
	}

	err = archive(xsources, basePath, outputPath)
	if err != nil {
		st.Step(terminal.StatusError, "Archive failed")
		return nil, err
	}

	st.Step(terminal.StatusOK, "Archive saved to '"+outputPath+"'")

	return &Archive{OutputPath: outputPath}, nil
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// expandSource returns the list of files recursively under a path.
func expandSource(source string, ignoreList []string) ([]string, error) {
	ignoreLookup := make(map[string]struct{})
	for _, i := range ignoreList {
		ignoreLookup[filepath.Clean(i)] = struct{}{}
	}

	var sources []string

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		_, ignore := ignoreLookup[rel]

		if info.IsDir() {
			if ignore {
				return filepath.SkipDir
			}

			return nil
		}

		if !ignore {
			sources = append(sources, path)
		}

		return nil
	}

	err := filepath.Walk(source, walker)
	if err != nil {
		return nil, err
	}

	return sources, nil
}
