package jsonutils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/momentag/ebliss/sdk/helpers/compressutils"
)

func EncodeJSON(in interface{}) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("input is nil")
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(in); err != nil {
		return nil, fmt.Errorf("error encoding input - %v", err.Error())
	}

	return buf.Bytes(), nil
}

func EncodeJSONAndCompress(in interface{}, config *compressutils.CompressionConfig) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("input is nil")
	}

	encoded, err := EncodeJSON(in)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = &compressutils.CompressionConfig{
			Type:                 compressutils.CompressionTypeGzip,
			GzipCompressionLevel: gzip.BestCompression,
		}
	}
	return compressutils.Compress(encoded, config)
}

func DecodeJSON(data []byte, out interface{}) error {
	if data == nil || len(data) == 0 {
		return fmt.Errorf("data being decoded is nil")
	}
	if out == nil {
		return fmt.Errorf("out parameter is nil")
	}

	decompressedBytes, uncompressed, err := compressutils.Decompress(data)
	if err != nil {
		return fmt.Errorf("failed to decompress jsonutils - %v", err.Error())
	}

	if !uncompressed && (decompressedBytes == nil || len(decompressedBytes) == 0) {
		return fmt.Errorf("decompressed data that is being decoded is invalid")
	}

	if !uncompressed {
		data = decompressedBytes
	}

	return DecodeJsonFromReader(bytes.NewReader(data), out)
}

func DecodeJsonFromReader(r io.Reader, out interface{}) error {
	if r == nil {
		return fmt.Errorf("reader is nil")
	}
	if out == nil {
		return fmt.Errorf("output parameter is nil")
	}
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(out)
}
