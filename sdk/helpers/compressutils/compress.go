package compressutils

import (
	"bytes"
	"compress/gzip"
	"compress/lzw"
	"fmt"
	"io"

	"github.com/golang/snappy"
	"github.com/pierrec/lz4"
	"github.com/rs/zerolog/log"
)

const (
	CompressionTypeGzip              = "gzip"
	CompressionTypeGzipCanary   byte = 'G'
	CompressionTypeLZW               = "lzw"
	CompressionTypeLZWCanary    byte = 'L'
	CompressionTypeSnappy            = "snappy"
	CompressionTypeSnappyCanary byte = 's'
	CompressionTypeLZ4               = "lz4"
	CompressionTypeLZ4Canary    byte = '4'
)

type CompressionReadCloser struct {
	io.Reader
}

func (c *CompressionReadCloser) Close() error {
	return nil
}

type CompressionConfig struct {
	Type                 string
	GzipCompressionLevel int
}

func Compress(data []byte, config *CompressionConfig) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	var writer io.WriteCloser

	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	switch config.Type {
	case CompressionTypeLZW:
		buf.Write([]byte{CompressionTypeLZWCanary})
		writer = lzw.NewWriter(&buf, lzw.LSB, 8)
	case CompressionTypeLZ4:
		buf.Write([]byte{CompressionTypeLZ4Canary})
		writer = lz4.NewWriter(&buf)
	case CompressionTypeSnappy:
		buf.Write([]byte{CompressionTypeSnappyCanary})
		writer = snappy.NewBufferedWriter(&buf)
	case CompressionTypeGzip:
		buf.Write([]byte{CompressionTypeGzipCanary})
		switch {
		case config.GzipCompressionLevel == gzip.BestCompression,
			config.GzipCompressionLevel == gzip.BestSpeed,
			config.GzipCompressionLevel == gzip.DefaultCompression:
		default:
			config.GzipCompressionLevel = gzip.DefaultCompression
		}
		writer, err = gzip.NewWriterLevel(&buf, config.GzipCompressionLevel)
	default:
		return nil, fmt.Errorf("unsupported compression algorithm")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create a compression writer: %v", err.Error())
	}

	if writer == nil {
		return nil, fmt.Errorf("failed to create a compression writer")
	}

	if _, err = writer.Write(data); err != nil {
		return nil, fmt.Errorf("failed to compressutils data: %v", err.Error())
	}

	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err.Error())
	}

	return buf.Bytes(), nil
}

func Decompress(data []byte) ([]byte, bool, error) {
	var err error
	var reader io.ReadCloser
	if data == nil || len(data) == 0 {
		return nil, false, fmt.Errorf("data is empty")
	}

	canary := data[0]
	content := data[1:]

	check := func(arg []byte) (io.Reader, error) {
		if len(arg) < 2 {
			return nil, fmt.Errorf("data after canary is invalid")
		} else {
			return bytes.NewReader(arg), nil
		}
	}

	decompress := func(arg *io.ReadCloser) ([]byte, bool, error) {
		if reader == nil {
			return nil, false, fmt.Errorf("failed to create compression reader")
		}
		defer (*arg).Close()
		var buf bytes.Buffer
		if written, err := io.Copy(&buf, *arg); err != nil {
			return nil, false, fmt.Errorf("could not copy bytes - %v", err.Error())
		} else {
			log.Trace().Int64("bytes_written", written)
			return buf.Bytes(), false, nil
		}
	}

	switch canary {
	case CompressionTypeLZWCanary:
		if rdr, err := check(content); err == nil {
			reader = lzw.NewReader(rdr, lzw.LSB, 8)
			return decompress(&reader)
		}
	case CompressionTypeLZ4Canary:
		if rdr, err := check(content); err == nil {
			reader = &CompressionReadCloser{Reader: lz4.NewReader(rdr)}
			return decompress(&reader)
		}
	case CompressionTypeSnappyCanary:
		if rdr, err := check(content); err == nil {
			reader = &CompressionReadCloser{Reader: snappy.NewReader(rdr)}
			return decompress(&reader)
		}
	case CompressionTypeGzipCanary:
		if rdr, err := check(content); err == nil {
			reader, err = gzip.NewReader(rdr)
			return decompress(&reader)
		}
	default:
		return nil, true, nil
	}

	return nil, false, fmt.Errorf("failed to create compression reader - %v", err.Error())
}
