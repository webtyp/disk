package disk_test

import (
	"os"
	"path/filepath"
	"testing"

	"webtyp.com/disk"
	"webtyp.com/files"
	"webtyp.com/files/conformance"
)

func TestDiskConformance(t *testing.T) {
	conformance.Run(t, conformance.Factory{
		Name: "disk",
		New:  func(t *testing.T) files.ReadWriter { return disk.Files{Root: t.TempDir()} },
	})
}

// WriteFile never leaves a temporary file behind and never exposes a truncated file.
func TestWriteFile_AtomicLeavesNoTempAndKeeps0644(t *testing.T) {
	dir := t.TempDir()
	d := disk.Files{Root: dir}
	if err := d.WriteFile(".env", []byte("A=1\nB=2\n")); err != nil {
		t.Fatal(err)
	}
	if err := d.WriteFile(".env", []byte("A=9\nB=2\nC=3\n")); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("want only .env in the directory, got %d entries", len(entries))
	}
	info, _ := os.Stat(filepath.Join(dir, ".env"))
	if info.Mode().Perm() != 0644 {
		t.Fatalf("mode = %v, want 0644", info.Mode().Perm())
	}
}

// Without Root, paths are used as given.
func TestFiles_ZeroValueUsesPathsAsGiven(t *testing.T) {
	p := filepath.Join(t.TempDir(), "x.txt")
	if err := (disk.Files{}).WriteFile(p, []byte("x")); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(p); string(got) != "x" {
		t.Fatalf("got %q", got)
	}
}
