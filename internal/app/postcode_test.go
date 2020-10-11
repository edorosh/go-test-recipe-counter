package app

import "testing"

func TestPostcodeFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Postcode
		wantErr bool
	}{
		{
			"Fail on invalid string",
			args{"blabla"},
			Postcode(0),
			true,
		},
		{
			"Fail on 0",
			args{"0"},
			Postcode(0),
			true,
		},
		{
			"Fail on negative",
			args{"-10"},
			Postcode(0),
			true,
		},
		{
			"Pass on valid",
			args{"9999999999"},
			Postcode(9999999999),
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := PostcodeFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostcodeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PostcodeFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
