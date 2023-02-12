package nonjson

import "testing"

func TestBruteForceProductKeyMACAddress(t *testing.T) {
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
			name: "valid",
			args: args{
				encodedProductKey: "AAYAAAAAAAAAAAAAAAAAAExLCU/N0RxxvG7ZACnE9iyfm1zRK6acy5rtKA01mFtnuCkFSJQtmsmoAN7KVyfxVbUpwPvJNKc2tkQbezXSbITnPSKnp8i9uG+C8DB+9oISsuTL8L0v07TOOsAnrSq4fR4mAhwANTYsmoLYmpqhVDLH/VVisfqVFSZu72vTDf2rjYESalQNawzIH8qjEhS2dzDUTm4RWf122JiPTSccbg2V8b4XXLRSefvc4ctVmCVvmrWRX+Aosgn9z0VS5V1ABhitiDjBd4NK34wOoGtn0vTwdiAfjMH95U5Q+c4hCjWsUnTrlUrdH5OQgtCDGi7Nag==",
			},
			want:    "3cecef123456",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BruteForceProductKeyMACAddress(tt.args.encodedProductKey)
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
