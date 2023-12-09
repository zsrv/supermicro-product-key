package oob

import (
	"testing"

	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

func TestEncodeOOBProductKey(t *testing.T) {
	type args struct {
		macAddress netinternal.HardwareAddr
	}
	tests := []struct {
		name    string
		args    args
		want    ProductKey
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				macAddress: netinternal.HardwareAddr{0x3C, 0xEC, 0xEF, 0x12, 0x34, 0x56},
			},
			want:    ProductKey{0xCE, 0x27, 0xF6, 0x41, 0x9B, 0x04, 0x6B, 0x24, 0x5D, 0x04, 0x5D, 0x32},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeOOBProductKey(tt.args.macAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeOOBProductKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("EncodeOOBProductKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
