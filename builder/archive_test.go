package builder

import (
	"archive/zip"
	"reflect"
	"testing"
)

func zipFiles(path string) ([]string, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	defer zr.Close()

	var files []string
	for _, f := range zr.File {
		files = append(files, f.Name)
	}

	return files, nil
}

func Test_archive(t *testing.T) {
	type args struct {
		sources  []string
		basePath string
		ignore   []string
	}

	tests := map[string]struct {
		args    args
		want    []string
		wantErr bool
	}{
		"nested-dirs no parent folder": {
			args: args{
				sources:  []string{"testdata/nested-dirs", "testdata/only-files"},
				basePath: "testdata",
				ignore:   []string{"ignore/", "file.txt", "file1.txt"},
			},
			want: []string{
				"nested-dirs/dir1/file.txt",
				"nested-dirs/dir1/nested-dir1/file.txt",
				"nested-dirs/dir1/nested-dir1/nested-nested-dir1/file1.txt",
				"nested-dirs/dir1/nested-dir1/nested-nested-dir1/file2.txt",
				"nested-dirs/dir1/nested-dir1/nested-nested-dir1/file3.txt",
				"nested-dirs/dir2/file1.txt",
				"nested-dirs/dir2/file2.txt",
				"nested-dirs/nested-dirs.zip",
				"only-files/file2.txt",
				"only-files/file3.txt",
				"only-files/file4.txt",
			},
			wantErr: false,
		},
		"nested-dirs with parent folder": {
			args: args{
				sources:  []string{"testdata/nested-dirs"},
				basePath: "testdata/nested-dirs",
				ignore:   []string{"ignore/", "file.txt", "file1.txt"},
			},
			want: []string{
				"dir1/file.txt",
				"dir1/nested-dir1/file.txt",
				"dir1/nested-dir1/nested-nested-dir1/file1.txt",
				"dir1/nested-dir1/nested-nested-dir1/file2.txt",
				"dir1/nested-dir1/nested-nested-dir1/file3.txt",
				"dir2/file1.txt",
				"dir2/file2.txt",
				"nested-dirs.zip",
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			outputPath := t.TempDir() + "/" + name + "output.zip"
			var xsources []string

			for _, src := range tt.args.sources {
				xsrc, _ := expandSource(src, tt.args.ignore)
				xsources = append(xsources, xsrc...)
			}

			if err := archive(xsources, tt.args.basePath, outputPath); (err != nil) != tt.wantErr {
				t.Errorf("archive() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := zipFiles(outputPath)
			if err != nil {
				t.Fatalf("could not open zip file: %s", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("archive() got = %v, want %v", got, tt.want)
			}
		})
	}
}
