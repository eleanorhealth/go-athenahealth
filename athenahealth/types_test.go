package athenahealth

import "testing"

func TestNumberString_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		n       *NumberString
		args    args
		wantErr bool
	}{
		{
			name:    "string value",
			n:       new(NumberString),
			args:    args{data: []byte(`"55.01"`)},
			wantErr: false,
		},
		{
			name:    "negative string value",
			n:       new(NumberString),
			args:    args{data: []byte(`"-55.01"`)},
			wantErr: false,
		},
		{
			name:    "int value",
			n:       new(NumberString),
			args:    args{data: []byte(`55`)},
			wantErr: false,
		},
		{
			name:    "float64 value",
			n:       new(NumberString),
			args:    args{data: []byte(`55.01`)},
			wantErr: false,
		},
		{
			name:    "invalid bool value",
			n:       new(NumberString),
			args:    args{data: []byte(`false`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Balance.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
