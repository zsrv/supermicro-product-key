package oob

import "testing"

func TestBruteForceProductKeyMACAddress(t *testing.T) {
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
			name: "valid",
			args: args{
				productKey: "CE27-F641-9B04-6B24-5D04-5D32",
			},
			want:    "3cecef123456",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BruteForceProductKeyMACAddress(tt.args.productKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("BruteForceProductKeyMACAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BruteForceProductKeyMACAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
