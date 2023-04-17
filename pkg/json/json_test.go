package json

import "testing"

var validOOBProductKey = `{"ProductKey":{"Node":{"LicenseID":"1","LicenseName":"SFT-OOB-LIC","CreateDate":"20200921"},"Signature":"OAaLKLy5IEK9WnIdnyA9ew89qTKQrm1eu+Q84CbwjR7XG7JGYccec+3vS3y/kQRRej3DcNVQPWsasX86ROTT+LZFsNY2mIEbQ6+Y/Tmv6+jwYgbQjEN6CjI7ahyKcebN12+3cLvPZyRf3kDqgtcpfuw3Qeg8BbhhyHQk29yNp+NG0XbKn02sHTrskvAGgG0GGlDCT5YmNa0gDSMzsvt/eH9nskb5opQNE3j7MAMXbjpI7xVHRbmB2N5iSu8gQUj0/pmk615ztM/uB54ur3GninJRU74S9Kotz+JunJg4pprGyQW544ggmzklmtr3zCA3GK/d929eZsVk5p8UxXG7wQ=="}}`
var validOOBProductKeyMACAddress = "ac1f6b3ddaec"

func TestVerifyProductKeySignature(t *testing.T) {
	type args struct {
		productKey string
		macAddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				productKey: validOOBProductKey,
				macAddress: validOOBProductKeyMACAddress,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyProductKeySignature(tt.args.productKey, tt.args.macAddress); (err != nil) != tt.wantErr {
				t.Errorf("VerifyProductKeySignature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
