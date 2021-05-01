package gobutils

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func EncodeInterface(in interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(in); err != nil {
		return buf, fmt.Errorf("cannot encode input (%v)", err.Error())
	} else {
		return buf, nil
	}
}

func DecodeData(data []byte, out interface{}) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("could not decode data (%v)", err.Error())
	} else {
		return nil
	}
}
