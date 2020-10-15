package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {

	startServer()

	raw := `{"ctime":30,"rtime":30,"certs":["client_test_crt.txt", "client_test_key.txt"],"ca":"client_test_crs.txt","proxy":"","keepalive":true,"trace":true}`
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "1", args: args{raw: raw}, want: "success", wantErr: false},
		{name: "2", args: args{raw: ""}, want: "success", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := fmt.Sprintf("http://192.168.5.115:9091/client?raw=%s", tt.args.raw)
			u, _ := url.Parse(req)
			q := u.Query()
			u.RawQuery = q.Encode()
			resp, err := http.Get(u.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp.StatusCode != 200 {
				t.Errorf("Request err,status:%d", resp.StatusCode)
				return
			}

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if !reflect.DeepEqual(string(body), tt.want) {
				t.Errorf("NewClient() = %v, want %v", string(body), tt.want)
			}
		})
	}
}