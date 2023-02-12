package oob

import "testing"

func TestEncodeOOBProductKey(t *testing.T) {
	type args struct {
		macAddress string
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
				macAddress: "3cecef123456",
			},
			want:    "CE27-F641-9B04-6B24-5D04-5D32",
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
			if got != tt.want {
				t.Errorf("EncodeOOBProductKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeProductKey(t *testing.T) {
	normalizedKey := "CE27-F641-9B04-6B24-5D04-5D32"

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normalized",
			args: args{
				key: normalizedKey,
			},
			want:    normalizedKey,
			wantErr: false,
		},
		{
			name: "uppercase without dashes",
			args: args{
				key: "CE27F6419B046B245D045D32",
			},
			want:    normalizedKey,
			wantErr: false,
		},
		{
			name: "lowercase without dashes",
			args: args{
				key: "ce27f6419b046b245d045d32",
			},
			want:    normalizedKey,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeProductKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeProductKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeProductKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeMACAddress(t *testing.T) {
	normalizedMACAddress := "3cecef123456"

	type args struct {
		macAddress string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normalized",
			args: args{
				macAddress: normalizedMACAddress,
			},
			want:    normalizedMACAddress,
			wantErr: false,
		},
		{
			name: "uppercase",
			args: args{
				macAddress: "3CECEF123456",
			},
			want:    normalizedMACAddress,
			wantErr: false,
		},
		{
			name: "uppercase with separators",
			args: args{
				macAddress: "3C:EC:EF:12:34:56",
			},
			want:    normalizedMACAddress,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeMACAddress(tt.args.macAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeMACAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeMACAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
