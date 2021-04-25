package jsonutils

import (
	"compress/gzip"
	"testing"

	"github.com/momentag/ebliss/sdk/helpers/compressutils"
	"github.com/stretchr/testify/assert"
)

var (
	expected = map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}
)

func TestEncodeJSON(t *testing.T) {
	t.Parallel()
	truth := []byte(`{"test":"data","validation":"process"}`)

	encodedJson, err := EncodeJSON(expected)
	assert.Nil(t, err)
	assert.NotNil(t, encodedJson)

	// test needs to remove the newline that the encoder adds
	assert.Equal(t, truth, encodedJson[:len(encodedJson)-1])

	_, err = EncodeJSON(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "input is nil")
}

func TestEncodeJSONAndCompress(t *testing.T) {
	t.Parallel()

	// Legit way
	compressedBytes, err := EncodeJSONAndCompress(expected, &compressutils.CompressionConfig{
		Type:                 compressutils.CompressionTypeGzip,
		GzipCompressionLevel: gzip.BestCompression,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, compressedBytes)
	assert.Equal(t, compressutils.CompressionTypeGzipCanary, compressedBytes[0])

	// check with nil config (by default uses gzip best compression)
	compressedBytes, err = EncodeJSONAndCompress(expected, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, compressedBytes)
	assert.Equal(t, compressutils.CompressionTypeGzipCanary, compressedBytes[0])

	// check with nil input
	_, err = EncodeJSONAndCompress(nil, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "input is nil")
}

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	// check legit
	input := []byte(`{"test":"data","validation":"process","other_object":{"boolean": true}}"`)
	expected := map[string]interface{}{
		"test":         "data",
		"validation":   "process",
		"other_object": map[string]interface{}{"boolean": true},
	}
	var actual map[string]interface{}
	err := DecodeJSON(input, &actual)
	assert.Nil(t, err)
	assert.Equal(t, actual, expected)

	// check nil input
	err = DecodeJSON(nil, &actual)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "data being decoded is nil")

	// check nil output
	err = DecodeJSON(input, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "out parameter is nil")

	// check compressed
	compressed, _ := EncodeJSONAndCompress(input, &compressutils.CompressionConfig{Type: compressutils.CompressionTypeGzip})
	decompressed := DecodeJSON(compressed, &actual)
	assert.NotNil(t, decompressed)
	assert.Equal(t, expected, actual)
}
