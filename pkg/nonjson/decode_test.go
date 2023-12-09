package nonjson

import (
	"reflect"
	"testing"
	"time"

	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

func TestDecodeProductKey(t *testing.T) {
	type args struct {
		encodedProductKey string
		macAddress        netinternal.HardwareAddr
	}
	tests := []struct {
		name    string
		args    args
		want    *ProductKey
		wantErr bool
	}{
		{
			name: "valid product key",
			args: args{
				encodedProductKey: validEncodedProductKey,
				macAddress:        validEncodedProductKeyMACAddress,
			},
			want: &ProductKey{
				FormatVersion:      0,
				SoftwareIdentifier: *SoftwareIdentifiers.ALL,
				SoftwareVersion:    "none",
				InvoiceNumber:      "none",
				CreationDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				ExpirationDate:     time.Unix(0, 0).UTC(),
				Property:           nil,
				SecretData:         []byte("76aa41c893a4e32a34203d8d47b0faab"),
				Checksum:           171,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEncodedProductKey(tt.args.encodedProductKey, tt.args.macAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEncodedProductKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseEncodedProductKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
