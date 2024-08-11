package main

import (
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"FOO":   "foo_value",
		"BAR":   "bar_value",
		"HELLO": "hello\x00world",
		"EMPTY": "",
		"TRIM":  "trim_value\t \n",
	}

	for name, content := range files {
		if err := os.WriteFile(dir+"/"+name, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err)
	}

	expected := map[string]EnvValue{
		"FOO":   {Value: "foo_value", NeedRemove: false},
		"BAR":   {Value: "bar_value", NeedRemove: false},
		"HELLO": {Value: "hello\nworld", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: true},
		"TRIM":  {Value: "trim_value", NeedRemove: false},
	}

	for key, expectedValue := range expected {
		if val, ok := env[key]; !ok {
			t.Errorf("Expected key %s to be present", key)
		} else if val != expectedValue {
			t.Errorf("For key %s, expected %+v, got %+v", key, expectedValue, val)
		}
	}

	if len(env) != len(expected) {
		t.Errorf("Expected %d keys, but got %d", len(expected), len(env))
	}
}

func TestReadDirWithInvalidPath(t *testing.T) {
	_, err := ReadDir("/non/existent/path")
	if err == nil {
		t.Fatal("Expected an error for a non-existent path")
	}
}

func TestReadDirWithSubdirectories(t *testing.T) {
	dir := t.TempDir()

	subdir := dir + "/subdir"
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	files := map[string]string{
		"FOO": "foo_value",
		"BAR": "bar_value",
	}

	for name, content := range files {
		if err := os.WriteFile(dir+"/"+name, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	if err := os.WriteFile(subdir+"/IGNORED", []byte("ignored_value"), 0o644); err != nil {
		t.Fatalf("Failed to create test file in subdir: %v", err)
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err)
	}

	expected := map[string]EnvValue{
		"FOO": {Value: "foo_value", NeedRemove: false},
		"BAR": {Value: "bar_value", NeedRemove: false},
	}

	for key, expectedValue := range expected {
		if val, ok := env[key]; !ok {
			t.Errorf("Expected key %s to be present", key)
		} else if val != expectedValue {
			t.Errorf("For key %s, expected %+v, got %+v", key, expectedValue, val)
		}
	}

	if len(env) != len(expected) {
		t.Errorf("Expected %d keys, but got %d", len(expected), len(env))
	}
}
