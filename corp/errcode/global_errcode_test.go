package errcode

import "testing"

func TestError(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "-1",
			args:    args{i: -1},
			wantErr: true,
		},
		{
			name:    "1",
			args:    args{i: 0},
			wantErr: false,
		},
		{
			name:    "-2",
			args:    args{i: -2},
			wantErr: true,
		},
		{
			name:    "40001",
			args:    args{i: 40001},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Error(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
