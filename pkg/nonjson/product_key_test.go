package nonjson

import (
	"testing"
	"time"

	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

func TestProductKey_Encode(t *testing.T) {
	type fields struct {
		FormatVersion      byte
		SoftwareIdentifier SoftwareIdentifier
		SoftwareVersion    string
		InvoiceNumber      string
		CreationDate       time.Time
		ExpirationDate     time.Time
		Property           []byte
		SecretData         []byte
		Checksum           byte
	}
	type args struct {
		macAddress netinternal.HardwareAddr
	}

	validProductKey := fields{
		FormatVersion:      0,
		SoftwareIdentifier: *SoftwareIdentifiers.ALL,
		SoftwareVersion:    "none",
		InvoiceNumber:      "none",
		CreationDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		ExpirationDate:     time.Unix(0, 0).UTC(),
		Property:           nil,
		SecretData:         []byte("76aa41c893a4e32a34203d8d47b0faab"),
		Checksum:           171,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "valid",
			fields: validProductKey,
			args: args{
				macAddress: validEncodedProductKeyMACAddress,
			},
			want:    validEncodedProductKey,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := &ProductKey{
				FormatVersion:      tt.fields.FormatVersion,
				SoftwareIdentifier: tt.fields.SoftwareIdentifier,
				SoftwareVersion:    tt.fields.SoftwareVersion,
				InvoiceNumber:      tt.fields.InvoiceNumber,
				CreationDate:       tt.fields.CreationDate,
				ExpirationDate:     tt.fields.ExpirationDate,
				Property:           tt.fields.Property,
				SecretData:         tt.fields.SecretData,
				Checksum:           tt.fields.Checksum,
			}
			got, err := pk.Encode(tt.args.macAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
