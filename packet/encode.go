package packet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
)

func Encode(r string) []byte {

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(r), &m)
	if err != nil {
		panic(err)
	}

	encoder := msgpack.GetEncoder()
	encoder.UseCompactInts(true)
	encoder.UseCompactFloats(true)
	buf := new(bytes.Buffer)
	encoder.Reset(buf)

	// 编码
	encode(encoder, m)

	//res := base64.StdEncoding.EncodeToString(buf.Bytes())
	return buf.Bytes()
}

func encode(encoder *msgpack.Encoder, r interface{}) {

	switch r.(type) {
	case map[string]interface{}:
		m := r.(map[string]interface{})
		fmt.Println("len :", len(m))
		encoder.EncodeMapLen(len(m))
		for k, v := range m {
			encoder.EncodeString(k)
			encode(encoder, v)
		}
	case string:
		encoder.EncodeString(r.(string))
	case []interface{}:
		arr := r.([]interface{})
		if len(arr) == 0 {
			return
		}
		encoder.EncodeArrayLen(len(arr))
		for _, v := range arr {
			encode(encoder, v)
		}
	case int, int8, int16, int32:
		encoder.EncodeInt32(r.(int32))
	case float32:
		encoder.EncodeInt32(int32(r.(float32)))
	case float64:
		encoder.EncodeInt32(int32(r.(float64)))
	case int64:
		encoder.EncodeInt64(r.(int64))
	default:
		encoder.Encode(r)
	}
}
