// Source code file, created by Developer@YAN_YING_SONG.

package xjs

import (
	"bytes"
	"encoding/json"

	"xie/go/errcause"

	"github.com/json-iterator/go"
)

var jsonFast = jsoniter.Config{EscapeHTML: false}.Froze()

func Unmarshal(b []byte, v interface{}) error {
	// Fast Unmarshal.

	if err := jsonFast.Unmarshal(b, v); err != nil {
		println(err.Error())
		if err = json.Unmarshal(b, v); err != nil {
			return err
		}
	}

	return nil
}

func ToJsonBytesFast(v interface{}) []byte {
	// Struct value to Json bytes, but escapeHTML = false.
	// Fast Marshal.

	b, err := jsonFast.Marshal(v)
	if err != nil {
		println(err.Error())
		return ToJsonBytes(v, false)
	}

	return b
}

func ToJsonBytes(v interface{}, format ...bool) []byte {
	// Struct value to Json bytes, but escapeHTML = false.

	buf := bytes.NewBuffer(make([]byte, 0, 128))
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if len(format) > 0 {
		if format[0] {
			enc.SetIndent("", "  ")
		}
	}
	if err := enc.Encode(v); err != nil {
		panic(err) // 编码失败则意味着致命的程序错误
	}
	data := buf.Bytes()
	n := len(data) - 1
	if data[n] == '\n' {
		return data[:n]
	}

	return data
}

func JsonBytesAppend(values ...[]byte) []byte {
	// Provide splicing of multiple sets of JSON data.

	l := 0
	for i := 0; i < len(values); i++ {
		if i != 0 {
			l += len(values[i-1]) - 1
			values[i][0] = ','
			values[0] = append(values[0][:l], values[i]...)
		}
	}

	return values[0]
}

func JsonBytesSetValue(m map[string]interface{}, value interface{}, keys ...string) error {
	// Modify values of arbitrary depth in JsonBytes. (Arrays are not supported)

	// Dive the data level.
	current := m
	for i := 0; i < len(keys)-1; i++ {
		v, ok := current[keys[i]]
		if !ok {
			// Create a new nested map if the key does not exist.
			next := make(map[string]interface{})
			current[keys[i]] = next
			current = next
		} else {
			current = v.(map[string]interface{})
		}
	}

	// Set value.
	name := keys[len(keys)-1]
	current[name] = value

	return nil
}

func InterfaceBind(value, bind interface{}) error {
	// Interface value are bind to struct.

	data := ToJsonBytes(value, false)
	if err := json.Unmarshal(data, bind); err != nil {
		return errcause.LinkErr(err, "json.Unmarshal")
	}

	return nil
}
