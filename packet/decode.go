package packet

import (
	"encoding/base64"
	"encoding/json"
	"github.com/ugorji/go/codec"
)

//Decode
func Decode(p []byte) []byte {
	m := make(map[string]interface{})
	codec.NewDecoderBytes(p, new(codec.MsgpackHandle)).Decode(&m)
	tmp := mapHandler(m)

	bytes, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	return bytes
}

func mapHandler(param interface{}) map[string]interface{} {
	tmp := make(map[string]interface{})
	switch param.(type) {
	case nil:
		return tmp

	case map[string]interface{}:
		for k, v := range param.(map[string]interface{}) {
			switch v.(type) {
			case map[interface{}]interface{}:
				tmp[k] = mapHandler(v)
			case []uint8:
				tmp[k] = B2S(v.([]uint8))

			default:
				tmp[k] = v
			}
		}
	case map[interface{}]interface{}:
		for k, v := range param.(map[interface{}]interface{}) {
			switch k.(type) {
			case string:
				switch v.(type) {
				case map[interface{}]interface{}:
					tmp[k.(string)] = mapHandler(v)
					continue
				case []interface{}:
					tmp[k.(string)] = B2SinArraySlice(v.([]interface{}))
				default:
					tmp[k.(string)] = v
				}
			default:
				continue
			}
		}
	}
	return tmp
}

func B2SinArraySlice(array []interface{}) []interface{} {
	tmp := make([]interface{}, 0)
	for _, v := range array {
		switch v.(type) {
		case []interface{}:
			tmp = append(tmp, B2SinArraySlice(v.([]interface{})))
		case []uint8:
			tmp = append(tmp, B2S(v.([]uint8)))
		default:
			tmp = append(tmp, v)

		}
	}
	return tmp
}

func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func DecodeString(p string) []byte {
	decodeBytes, err := base64.StdEncoding.DecodeString(p)
	if err != nil {
		panic(err)
	}
	return Decode(decodeBytes)
}
