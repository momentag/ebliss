package compressutils

import (
	"compress/gzip"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress_CompressDecompress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		compressionType   string
		compressionConfig CompressionConfig
		canary            byte
	}{
		{
			"GZIP Default Implicit",
			CompressionConfig{Type: CompressionTypeGzip},
			CompressionTypeGzipCanary,
		},
		{
			"GZIP Default Explicit",
			CompressionConfig{Type: CompressionTypeGzip, GzipCompressionLevel: gzip.DefaultCompression},
			CompressionTypeGzipCanary,
		},
		{
			"GZIP Best Speed",
			CompressionConfig{Type: CompressionTypeGzip, GzipCompressionLevel: gzip.BestSpeed},
			CompressionTypeGzipCanary,
		},
		{
			"GZIP Best Compression",
			CompressionConfig{Type: CompressionTypeGzip, GzipCompressionLevel: gzip.BestCompression},
			CompressionTypeGzipCanary,
		},
		{
			"Snappy",
			CompressionConfig{Type: CompressionTypeSnappy},
			CompressionTypeSnappyCanary,
		},
		{
			"LZW",
			CompressionConfig{Type: CompressionTypeLZW},
			CompressionTypeLZWCanary,
		},
		{
			"LZ4",
			CompressionConfig{Type: CompressionTypeLZ4},
			CompressionTypeLZ4Canary,
		},
	}

	inputBytes := []byte(`{"sample":"data","verification":"process"}`)

	for _, test := range tests {
		compressedBytes, err := Compress(inputBytes, &test.compressionConfig)
		assert.Nil(t, err)
		assert.NotNil(t, compressedBytes)
		assert.NotEmpty(t, compressedBytes)
		assert.Contains(t, compressedBytes, test.canary)
		assert.Equal(t, test.canary, compressedBytes[0])

		decompressedBytes, wasNotCompressed, err := Decompress(compressedBytes)
		assert.Nil(t, err)
		assert.False(t, wasNotCompressed)
		assert.NotNil(t, decompressedBytes)
		assert.NotEmpty(t, decompressedBytes)
		assert.Equal(t, inputBytes, decompressedBytes)

		// assume data was not compressed by removing canary byte
		pseudoCompressed := compressedBytes[1:]
		_, flag, err := Decompress(pseudoCompressed)
		assert.Nil(t, err)
		assert.True(t, flag)
	}
}

func TestInvalidConfigs(t *testing.T) {
	t.Parallel()
	inputBytes := []byte(`{"sample":"data","verification":"process"}`)

	_, err := Compress(inputBytes, nil)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("config is nil"), err)

	_, err = Compress(inputBytes, &CompressionConfig{})
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("unsupported compression algorithm"), err)
}

func TestWithEmptyData(t *testing.T) {
	compressedBytes := []byte{}
	_, _, err := Decompress(compressedBytes)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("data is empty"), err)
}
