// Copyright Kirk Rader 2024

package utilities

import (
	"archive/zip"
	"io"
	"os"
)

// Type of function used to process each entry in a zip archive.
type ZipHandler func(entry *zip.File)

// Apply the given handler to each entry in the given zip file.
func ForEachZipEntry(handler ZipHandler, archive io.ReaderAt, size int64) error {
	zip, err := zip.NewReader(archive, size)
	if err != nil {
		return err
	}
	for _, entry := range zip.File {
		handler(entry)
	}
	return nil
}

// Apply the given handler to each entry in the given zip archive.
func ForZipFile(handler ZipHandler, file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	return ForEachZipEntry(handler, file, info.Size())
}

// Apply the given handler to each entry in the given zip archive.
func ForZipReader(handler ZipHandler, reader io.Reader) error {
	file, err := os.CreateTemp(os.TempDir(), "ForZipReader")
	if err != nil {
		return err
	}
	// deferred functions are invoked in reverse order
	defer os.Remove(file.Name()) // invoked second
	defer file.Close()           // invoked first
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	return ForZipFile(handler, file)
}
