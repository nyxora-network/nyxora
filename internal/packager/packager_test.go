package packager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewPackager(t *testing.T) {
	dir := t.TempDir()
	p := NewPackager(dir)
	if p == nil {
		t.Fatal("NewPackager returned nil")
	}
}

func TestPackagerPackUnpack(t *testing.T) {
	dir := t.TempDir()
	p := NewPackager(dir)

	// Create source dir with a file
	srcDir := filepath.Join(dir, "source")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "test.txt"), []byte("hello world"), 0644)

	// Pack
	err := p.Pack("test", srcDir)
	if err != nil {
		t.Fatalf("Pack failed: %v", err)
	}

	// Check archive exists
	if !p.Exists("test") {
		t.Error("archive should exist after Pack")
	}

	// Unpack
	destDir := filepath.Join(dir, "dest")
	err = p.Unpack("test", destDir)
	if err != nil {
		t.Fatalf("Unpack failed: %v", err)
	}

	// Verify content
	data, err := os.ReadFile(filepath.Join(destDir, "test.txt"))
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("content should be 'hello world', got '%s'", string(data))
	}
}

func TestPackagerList(t *testing.T) {
	dir := t.TempDir()
	p := NewPackager(dir)

	list, err := p.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("empty packager should have 0 items, got %d", len(list))
	}
}

func TestPackagerRemove(t *testing.T) {
	dir := t.TempDir()
	p := NewPackager(dir)

	srcDir := filepath.Join(dir, "source")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "test.txt"), []byte("hello"), 0644)

	p.Pack("test", srcDir)

	if !p.Exists("test") {
		t.Error("archive should exist before Remove")
	}

	err := p.Remove("test")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if p.Exists("test") {
		t.Error("archive should not exist after Remove")
	}
}

func TestPackagerUnpackNonexistent(t *testing.T) {
	dir := t.TempDir()
	p := NewPackager(dir)

	err := p.Unpack("nonexistent", filepath.Join(dir, "dest"))
	if err == nil {
		t.Error("Unpack of nonexistent archive should fail")
	}
}

func TestPackagerPackNonexistent(t *testing.T) {
	dir := t.TempDir()
	p := NewPackager(dir)

	err := p.Pack("test", "/nonexistent/path")
	if err == nil {
		t.Error("Pack of nonexistent source should fail")
	}
}
