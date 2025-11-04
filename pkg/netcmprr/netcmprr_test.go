package netcmprr

import "testing"

func TestIsTrustedSubnet(t *testing.T) {
	type args struct {
		trustedSubnet string
		requestIP     string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"positive case",
			args{"192.168.1.10/24", "192.168.1.20"},
			true,
			false,
		},
		{
			"negative case",
			args{"127.0.0.1/24", "192.168.1.20"},
			false,
			false,
		},
		{
			"corner case",
			args{"", ""},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsTrustedSubnet(tt.args.trustedSubnet, tt.args.requestIP)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsTrustedSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsTrustedSubnet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
