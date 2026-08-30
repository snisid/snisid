package main

import (
	"testing"
)

func TestGetEnv_Default(t *testing.T) {
	if v := getEnv("NONEXISTENT_TRANSPORT_VAR", "fallback"); v != "fallback" {
		t.Errorf("expected fallback, got %s", v)
	}
}

func TestGetEnv_FromEnv(t *testing.T) {
	t.Setenv("TRANSPORT_TEST_VAR", "from-env")
	if v := getEnv("TRANSPORT_TEST_VAR", "fallback"); v != "from-env" {
		t.Errorf("expected from-env, got %s", v)
	}
}
