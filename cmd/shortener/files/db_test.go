package files

import "testing"

func TestDirExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "default test", args: args{path: "tmp/short-url-db.json"}, wantErr: false},
		{name: "incorrect path test 1", args: args{path: "helloworld"}, wantErr: true},
		{name: "incorrect path test 2", args: args{path: "helloworld/"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DirExists(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("DirExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
