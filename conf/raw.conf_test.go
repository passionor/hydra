package conf

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/micro-plat/lib4go/security/md5"
)

func TestNewRawConfByMap(t *testing.T) {
	dataN := map[string]interface{}{}
	dataV := map[string]interface{}{"t": "x"}
	dataB, _ := json.Marshal(dataV)
	type args struct {
		data    map[string]interface{}
		version int32
	}
	tests := []struct {
		name    string
		args    args
		wantC   *RawConf
		wantErr bool
	}{
		{name: "nil数据初始化", args: args{data: nil, version: 0}, wantC: &RawConf{data: nil, version: 0, raw: []byte("null"), signature: md5.EncryptBytes([]byte("null"))}, wantErr: false},
		{name: "空数据初始化", args: args{data: dataN, version: 0}, wantC: &RawConf{data: dataN, version: 0, raw: []byte("{}"), signature: md5.EncryptBytes([]byte("{}"))}, wantErr: false},
		{name: "对象数据初始化", args: args{data: dataV, version: 0}, wantC: &RawConf{data: dataV, version: 0, raw: dataB, signature: md5.EncryptBytes(dataB)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewRawConfByMap(tt.args.data, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRawConfByMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewRawConfByMap() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestNewRawConfByJson(t *testing.T) {
	type args struct {
		message []byte
		version int32
	}
	tests := []struct {
		name    string
		args    args
		wantC   *RawConf
		wantErr bool
	}{
		{name: "反序列化失败初始化", args: args{message: []byte("{sdfsdfsdf}"), version: 0}, wantC: &RawConf{data: nil, version: 0, raw: []byte("{sdfsdfsdf}"), signature: md5.EncryptBytes([]byte("{sdfsdfsdf}"))}, wantErr: true},
		{name: "正常初始化", args: args{message: []byte(`{"sss":"11"}`), version: 0}, wantC: &RawConf{data: map[string]interface{}{"sss": "11"}, version: 0, raw: []byte(`{"sss":"11"}`), signature: md5.EncryptBytes([]byte(`{"sss":"11"}`))}, wantErr: false},
		{name: "无data,正常初始化", args: args{message: []byte(`test1`), version: 0}, wantC: &RawConf{data: nil, version: 0, raw: []byte("test1"), signature: md5.EncryptBytes([]byte(`test1`))}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewRawConfByJson(tt.args.message, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRawConfByJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewRawConfByJson() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestRawConf_GetString(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	dn, _ := NewRawConfByMap(nil, 1)
	dataN := fields{data: dn.data, version: dn.version, raw: dn.raw, signature: dn.signature}
	dv, _ := NewRawConfByMap(map[string]interface{}{"test": "1"}, 1)
	dataV := fields{data: dv.data, version: dv.version, raw: dv.raw, signature: dv.signature}
	type args struct {
		key string
		def []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantR  string
	}{
		{name: "空对象获取", fields: dataN, args: args{key: "", def: []string{}}, wantR: ""},
		{name: "空对象获取", fields: dataN, args: args{}, wantR: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				data:      tt.fields.data,
			}
			if gotR := j.GetString(tt.args.key, tt.args.def...); gotR != tt.wantR {
				t.Errorf("RawConf.GetString() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestRawConf_GetInt(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		key string
		def []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				data:      tt.fields.data,
			}
			if got := j.GetInt(tt.args.key, tt.args.def...); got != tt.want {
				t.Errorf("RawConf.GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawConf_GetStrings(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		key string
		def []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantR  []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				data:      tt.fields.data,
			}
			if gotR := j.GetStrings(tt.args.key, tt.args.def...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("RawConf.GetStrings() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestRawConf_GetArray(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		key string
		def []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantR  []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				data:      tt.fields.data,
			}
			if gotR := j.GetArray(tt.args.key, tt.args.def...); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("RawConf.GetArray() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestRawConf_GetJSON(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		section string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantR       []byte
		wantVersion int32
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				data:      tt.fields.data,
			}
			gotR, gotVersion, err := j.GetJSON(tt.args.section)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawConf.GetJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("RawConf.GetJSON() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("RawConf.GetJSON() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

func TestRawConf_HasSection(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		section string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				data:      tt.fields.data,
			}
			if got := j.HasSection(tt.args.section); got != tt.want {
				t.Errorf("RawConf.HasSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawConf_GetSection(t *testing.T) {
	type fields struct {
		raw       json.RawMessage
		version   int32
		signature string
		data      map[string]interface{}
	}
	type args struct {
		section string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantC   *RawConf
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &RawConf{
				raw:       tt.fields.raw,
				version:   tt.fields.version,
				signature: tt.fields.signature,
				data:      tt.fields.data,
			}
			gotC, err := j.GetSection(tt.args.section)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawConf.GetSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("RawConf.GetSection() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func Test_parseBool(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := parseBool(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotValue != tt.wantValue {
				t.Errorf("parseBool() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}
