package builder

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// archive create a zip archive in the outputPath containing the source files.
// basePath controls the root folder inside the archive, for example given the source
// []string{"/home/parent/dir1/file1.txt", "/home/parent/dir1/file2.txt", "/home/parent/file.txt"},
// if the basePath is "/home", then the output archive will contain the parent directory.
// If the basePath is "/home/parent", then the archive will only have the content of "parent"
// but not the directory "parent" itself.
func archive(sources []string, basePath, outputPath string) error {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	for _, source := range sources {
		err := func() error {
			abs, err := filepath.Abs(source)
			if err != nil {
				return err
			}

			// make sure we are not trying to add the output archive
			// into the archive.
			if outputPath == abs {
				return nil
			}

			sourceFile, err := os.Open(source)
			if err != nil {
				return err
			}

			defer sourceFile.Close()

			rel, err := filepath.Rel(basePath, source)
			if err != nil {
				return err
			}

			f, err := w.Create(rel)
			if err != nil {
				return err
			}

			_, err = io.Copy(f, sourceFile)
			if err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
