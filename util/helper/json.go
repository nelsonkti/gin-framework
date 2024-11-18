package helper

import (
	jsoniter "github.com/json-iterator/go"
)

const JsonContentType = "application/json"

func Marshal(v interface{}) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
}

func UmMarshal(data []byte, v any) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, v)
}

func UnMarshalWithInterface(data interface{}, v any) error {
	d, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(data)
	if err != nil {
		return err
	}

	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(d, &v)
}
