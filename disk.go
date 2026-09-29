// Package disk is the implementation of the webtyp.com/files contract on the machine's file
// system, for servers and development tools. In a browser use webtyp.com/opfs instead.
package disk

import (
	"os"
	"path/filepath"

	"webtyp.com/files"
)

// Files reads and writes whole files on disk. The zero value works with paths as given; set
// Root to keep every path inside one directory.
type Files struct {
	Root string // when not empty, every path is joined to Root
}

var (
	_ files.ReadWriter = Files{}
	_ files.Appender   = Files{}
)

func (f Files) path(p string) string {
	if f.Root == "" {
		return p
	}
	return filepath.Join(f.Root, p)
}

// ReadFile returns the file's contents, or files.ErrNotExist when it does not exist.
func (f Files) ReadFile(p string) ([]byte, error) {
	data, err := os.ReadFile(f.path(p))
	if os.IsNotExist(err) {
		return nil, files.ErrNotExist
	}
	return data, err
}

// WriteFile replaces the file's contents atomically. The data is written to a temporary file in
// the same directory, which is then renamed over the target, so a reader never observes a
// truncated or half-written file (os.WriteFile truncates first and writes second, and a reader
// in that window cannot tell an emptied file from a genuinely empty one). The file is 0644.
func (f Files) WriteFile(p string, data []byte) error {
	return writeAtomic(f.path(p), data)
}

// AppendFile adds data to the end of the file, creating it (0644) if needed.
func (f Files) AppendFile(p string, data []byte) error {
	fh, err := os.OpenFile(f.path(p), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, werr := fh.Write(data)
	cerr := fh.Close()
	if werr != nil {
		return werr
	}
	return cerr
}

// writeAtomic is WriteFile's mechanism. The temporary file must live in the target's own
// directory: os.Rename is only atomic within one file system.
func writeAtomic(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	err = func() error {
		if _, err := tmp.Write(data); err != nil {
			return err
		}
		if err := tmp.Close(); err != nil {
			return err
		}
		// CreateTemp makes the file 0600; the file it replaces is readable by other tools.
		if err := os.Chmod(name, 0644); err != nil {
			return err
		}
		return os.Rename(name, path)
	}()
	if err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
	}
	return err
}
