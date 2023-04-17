package json

import (
	"testing"
)

func TestBruteForceMACAddressFromString(t *testing.T) {
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
				productKey: validOOBProductKey,
			},
			want:    validOOBProductKeyMACAddress,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BruteForceMACAddressFromString(tt.args.productKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("BruteForceMACAddressFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BruteForceMACAddressFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
