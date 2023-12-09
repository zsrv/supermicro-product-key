package json

import (
	"testing"

	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

var (
	validProductKey    = `{"ProductKey":{"Node":{"LicenseID":"1","LicenseName":"SFT-OOB-LIC","CreateDate":"20200921"},"Signature":"OAaLKLy5IEK9WnIdnyA9ew89qTKQrm1eu+Q84CbwjR7XG7JGYccec+3vS3y/kQRRej3DcNVQPWsasX86ROTT+LZFsNY2mIEbQ6+Y/Tmv6+jwYgbQjEN6CjI7ahyKcebN12+3cLvPZyRf3kDqgtcpfuw3Qeg8BbhhyHQk29yNp+NG0XbKn02sHTrskvAGgG0GGlDCT5YmNa0gDSMzsvt/eH9nskb5opQNE3j7MAMXbjpI7xVHRbmB2N5iSu8gQUj0/pmk615ztM/uB54ur3GninJRU74S9Kotz+JunJg4pprGyQW544ggmzklmtr3zCA3GK/d929eZsVk5p8UxXG7wQ=="}}`
	validProductKeyMAC = netinternal.HardwareAddr{0xAC, 0x1F, 0x6B, 0x3D, 0xDA, 0xEC}
	invalidProductKey  = `{"ProductKey":{"Node":{"LicenseID":"0","LicenseName":"INVALID-LIC","CreateDate":"00000000"},"Signature":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="}}`
)

func TestBruteForceMACAddressFromString(t *testing.T) {
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
				productKey: validProductKey,
			},
			want:    validProductKeyMAC.String(),
			wantErr: false,
		},
		{
			name: "invalid key",
			args: args{
				productKey: invalidProductKey,
			},
			wantErr: true,
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
