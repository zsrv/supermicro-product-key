package nonjson

import (
	"reflect"
	"testing"
	"time"
)

func TestDecodeProductKey(t *testing.T) {
	type args struct {
		encodedProductKey string
		macAddress        string
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
				encodedProductKey: "AAYAAAAAAAAAAAAAAAAAAExLCU/N0RxxvG7ZACnE9iyfm1zRK6acy5rtKA01mFtnuCkFSJQtmsmoAN7KVyfxVbUpwPvJNKc2tkQbezXSbITnPSKnp8i9uG+C8DB+9oISsuTL8L0v07TOOsAnrSq4fR4mAhwANTYsmoLYmpqhVDLH/VVisfqVFSZu72vTDf2rjYESalQNawzIH8qjEhS2dzDUTm4RWf122JiPTSccbg2V8b4XXLRSefvc4ctVmCVvmrWRX+Aosgn9z0VS5V1ABhitiDjBd4NK34wOoGtn0vTwdiAfjMH95U5Q+c4hCjWsUnTrlUrdH5OQgtCDGi7Nag==",
				macAddress:        "3cecef123456",
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
			got, err := DecodeProductKey(tt.args.encodedProductKey, tt.args.macAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeProductKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeProductKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
