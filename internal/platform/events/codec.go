package events

import (
	"encoding/json"
	"fmt"

	"github.com/hamba/avro/v2"
)

// Codec defines the interface for encoding and decoding events
type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// JSONCodec implements Codec for JSON
type JSONCodec struct{}

func (c *JSONCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *JSONCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// AvroCodec implements Codec for Avro
type AvroCodec struct {
	schema avro.Schema
}

func NewAvroCodec(schemaStr string) (*AvroCodec, error) {
	s, err := avro.Parse(schemaStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse avro schema: %w", err)
	}
	return &AvroCodec{schema: s}, nil
}

func (c *AvroCodec) Encode(v interface{}) ([]byte, error) {
	return avro.Marshal(c.schema, v)
}

func (c *AvroCodec) Decode(data []byte, v interface{}) error {
	return avro.Unmarshal(c.schema, data, v)
}
