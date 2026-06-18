package events

import (
	"testing"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestJSONCodec_EncodeDecode(t *testing.T) {
	codec := &JSONCode{}
	original := person{Name: "Jean", Age: 30}

	data, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("Encoded data should not be empty")
	}

	var decoded person
	err = codec.Decode(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded.Name != "Jean" {
		t.Errorf("Name = %s, want Jean", decoded.Name)
	}
	if decoded.Age != 30 {
		t.Errorf("Age = %d, want 30", decoded.Age)
	}
}

func TestJSONCodec_EncodeNil(t *testing.T) {
	codec := &JSONCodec{}
	_, err := codec.Encode(nil)
	if err != nil {
		t.Fatalf("Encode nil failed: %v", err)
	}
}

func TestJSONCodec_Decode_InvalidData(t *testing.T) {
	codec := &JSONCodec{}
	var result map[string]interface{}
	err := codec.Decode([]byte(`not json`), &result)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestJSONCodec_RoundTrip_Complex(t *testing.T) {
	codec := &JSONCodec{}
	input := map[string]interface{}{
		"identityId": "ID-001",
		"roles":      []string{"admin", "officer"},
		"metadata": map[string]interface{}{
			"ip": "10.0.0.1",
			"ts": 1234567890,
		},
		"score": 95.5,
	}

	data, err := codec.Encode(input)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var output map[string]interface{}
	err = codec.Decode(data, &output)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if output["identityId"] != "ID-001" {
		t.Errorf("identityId = %s, want ID-001", output["identityId"])
	}
}

func TestNewAvroCodec_InvalidSchema(t *testing.T) {
	_, err := NewAvroCodec("invalid schema")
	if err == nil {
		t.Error("Expected error for invalid Avro schema")
	}
}

func TestNewAvroCodec_ValidSchema(t *testing.T) {
	schema := `{
		"type": "record",
		"name": "TestEvent",
		"fields": [
			{"name": "id", "type": "string"},
			{"name": "value", "type": "int"}
		]
	}`
	codec, err := NewAvroCodec(schema)
	if err != nil {
		t.Fatalf("NewAvroCodec failed: %v", err)
	}
	if codec == nil {
		t.Fatal("NewAvroCodec returned nil")
	}
}

func TestJSONCodec_EncodeMap(t *testing.T) {
	codec := &JSONCodec{}
	data, err := codec.Encode(map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if string(data) != `{"key":"value"}` {
		t.Errorf("data = %s, want {\"key\":\"value\"}", string(data))
	}
}
