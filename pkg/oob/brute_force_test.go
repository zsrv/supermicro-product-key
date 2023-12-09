package oob

import "testing"

func TestBruteForceMACAddress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	type args struct {
		productKey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid key",
			args: args{
				productKey: "CE27-F641-9B04-6B24-5D04-5D32",
			},
			want:    "3CECEF123456",
			wantErr: false,
		},
		{
			name: "invalid key",
			args: args{
				productKey: "0000-0000-0000-0000-0000-0000",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BruteForceMACAddress(tt.args.productKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("BruteForceMACAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BruteForceMACAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
