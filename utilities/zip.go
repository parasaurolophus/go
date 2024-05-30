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
func ForZipFile(handler ZipHandler, archiveFile *os.File) error {
	info, err := archiveFile.Stat()
	if err != nil {
		return err
	}
	return ForEachZipEntry(handler, archiveFile, info.Size())
}

// Apply the given handler to the each entry in the given zip archive.
func ForZipReader(handler ZipHandler, archiveReader io.Reader) error {
	tempFile, err := os.CreateTemp(os.TempDir(), "ForZipReader")
	if err != nil {
		return err
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()
	_, err = io.Copy(tempFile, archiveReader)
	if err != nil {
		return err
	}
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return err
	}
	return ForZipFile(handler, tempFile)
}
