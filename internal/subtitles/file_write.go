package subtitles

import (
	"os"
	"path/filepath"
)

func writeFilePreserveMode(path string, data []byte) error {
	perm := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		perm = info.Mode().Perm()
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}

	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	_, err = tmp.Write(data)
	if err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	directory, err := os.Open(dir)
	if err != nil {
		return nil
	}
	defer func() {
		_ = directory.Close()
	}()

	if err := directory.Sync(); err != nil && !isSyncNotSupported(err) {
		return err
	}

	return nil
}

func isSyncNotSupported(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == "operation not supported" || err.Error() == "invalid argument"
}
