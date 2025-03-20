package registry

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/julien-fruteau/go-distribution-registry/internal/env"
)

var fsBlobsPath = env.GetEnvOrDefault("REG_BLOBS_PATH", "/var/lib/registry/docker/registry/v2/blobs/sha256")

func IsFileGzip(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	mb := ReadMagicBytes(f)
	return IsGzipMagicBytes(mb), nil
}

func ReadMagicBytes(r io.Reader) []byte {
	buf := make([]byte, 2)
	r.Read(buf)
	return buf
}

func IsGzipMagicBytes(b []byte) bool {
	return bytes.Equal(b, []byte{0x1F, 0x8B})
}

func WalkDirFnGzipBlobs(path string, d fs.DirEntry, err error, gzipBlobs *[]string) error {
	if err != nil {
		return err
	}
	if d.IsDir() || len(path) <= len(fsBlobsPath)+4 {
		return nil
	}

	gzip, err := IsFileGzip(path)
	if err != nil {
		return err
	}
	if gzip {
		*gzipBlobs = append(*gzipBlobs, path)
	}

	return nil
}

func WalkFs(root string) ([]string, error) {
	var gzipBlobs []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		return WalkDirFnGzipBlobs(path, d, err, &gzipBlobs)
	})
	if err != nil {
		return nil, err
	}

	return gzipBlobs, nil
}
