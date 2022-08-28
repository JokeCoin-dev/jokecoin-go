package utils

import (
	"github.com/vmihailenco/msgpack/v5"
)

/*
func Serialize(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}


func Deserialize(b []byte, v interface{}) error {
	dec := gob.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(v)
	return err
}
*/

func Serialize(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func Deserialize(b []byte, v interface{}) error {
	return msgpack.Unmarshal(b, v)
}

func MustSerialize(v interface{}) []byte {
	b, err := Serialize(v)
	PanicIfErr(err)
	return b
}

func MustDeserialize(b []byte, v interface{}) {
	err := Deserialize(b, v)
	PanicIfErr(err)
}
