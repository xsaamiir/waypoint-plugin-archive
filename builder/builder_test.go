package builder

import (
	"reflect"
	"testing"
)

func Test_expandSource(t *testing.T) {
	type args struct {
		path   string
		ignore []string
	}

	tests := []struct {
		args    args
		want    []string
		wantErr bool
	}{
		{
			args: args{
				path:   "testdata/nested-dirs",
				ignore: []string{"ignore/", "file.txt"},
			},
			want: []string{
				"testdata/nested-dirs/dir1/file.txt",
				"testdata/nested-dirs/dir1/nested-dir1/file.txt",
				"testdata/nested-dirs/dir1/nested-dir1/nested-nested-dir1/file1.txt",
				"testdata/nested-dirs/dir1/nested-dir1/nested-nested-dir1/file2.txt",
				"testdata/nested-dirs/dir1/nested-dir1/nested-nested-dir1/file3.txt",
				"testdata/nested-dirs/dir2/file1.txt",
				"testdata/nested-dirs/dir2/file2.txt",
				"testdata/nested-dirs/nested-dirs.zip",
			},
			wantErr: false,
		},
		{
			args:    args{path: "testdata/only-dirs"},
			want:    nil,
			wantErr: false,
		},
		{
			args: args{
				path:   "testdata/only-files",
				ignore: []string{"file1.txt"},
			},
			want: []string{
				"testdata/only-files/file2.txt",
				"testdata/only-files/file3.txt",
				"testdata/only-files/file4.txt",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.args.path, func(t *testing.T) {
			got, err := expandSource(tt.args.path, tt.args.ignore)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandSource() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expandSource() got = %v, want %v", got, tt.want)
			}
		})
	}
}
