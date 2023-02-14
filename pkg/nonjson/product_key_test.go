package nonjson

import (
	"testing"
	"time"
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
		macAddress string
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
				macAddress: "3cecef123456",
			},
			want:    "AAYAAAAAAAAAAAAAAAAAAExLCU/N0RxxvG7ZACnE9iyfm1zRK6acy5rtKA01mFtnuCkFSJQtmsmoAN7KVyfxVbUpwPvJNKc2tkQbezXSbITnPSKnp8i9uG+C8DB+9oISsuTL8L0v07TOOsAnrSq4fR4mAhwANTYsmoLYmpqhVDLH/VVisfqVFSZu72vTDf2rjYESalQNawzIH8qjEhS2dzDUTm4RWf122JiPTSccbg2V8b4XXLRSefvc4ctVmCVvmrWRX+Aosgn9z0VS5V1ABhitiDjBd4NK34wOoGtn0vTwdiAfjMH95U5Q+c4hCjWsUnTrlUrdH5OQgtCDGi7Nag==",
			wantErr: false,
		},
		{
			name:   "mac address with invalid character",
			fields: validProductKey,
			args: args{
				macAddress: "3cecef12345x",
			},
			wantErr: true,
		},
		{
			name:   "mac address too long",
			fields: validProductKey,
			args: args{
				macAddress: "3cecef1234567",
			},
			wantErr: true,
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
