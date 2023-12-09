package nonjson

import (
	"testing"

	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

var (
	validEncodedProductKey           = "AAYAAAAAAAAAAAAAAAAAAExLCU/N0RxxvG7ZACnE9iyfm1zRK6acy5rtKA01mFtnuCkFSJQtmsmoAN7KVyfxVbUpwPvJNKc2tkQbezXSbITnPSKnp8i9uG+C8DB+9oISsuTL8L0v07TOOsAnrSq4fR4mAhwANTYsmoLYmpqhVDLH/VVisfqVFSZu72vTDf2rjYESalQNawzIH8qjEhS2dzDUTm4RWf122JiPTSccbg2V8b4XXLRSefvc4ctVmCVvmrWRX+Aosgn9z0VS5V1ABhitiDjBd4NK34wOoGtn0vTwdiAfjMH95U5Q+c4hCjWsUnTrlUrdH5OQgtCDGi7Nag=="
	validEncodedProductKeyMACAddress = netinternal.HardwareAddr{0x3C, 0xEC, 0xEF, 0x12, 0x34, 0x56}
	invalidEncodedProductKey         = "AAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
)

func TestBruteForceMACAddress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	type args struct {
		encodedProductKey string
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
				encodedProductKey: validEncodedProductKey,
			},
			want:    validEncodedProductKeyMACAddress.String(),
			wantErr: false,
		},
		{
			name: "invalid key",
			args: args{
				encodedProductKey: invalidEncodedProductKey,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BruteForceMACAddress(tt.args.encodedProductKey)
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
