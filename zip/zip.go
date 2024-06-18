// Copyright Kirk Rader 2024

package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

// Type of function used to process each entry in a zip archive.
type ZipHandler func(*zip.File) error

// Apply the given handler to each entry in the given zip file. Terminate the
// loop upon first error.
func ForEachZipEntry(handler ZipHandler, readerAt io.ReaderAt, size int64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("ZipHandler panic: %v", r)
		}
	}()
	zipReader, err := zip.NewReader(readerAt, size)
	if err != nil {
		return
	}
	for _, entry := range zipReader.File {
		err = handler(entry)
		if err != nil {
			return
		}
	}
	return
}

// Apply the given handler to each entry in the given zip archive.
func ForEachZipEntryFromFile(handler ZipHandler, file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	return ForEachZipEntry(handler, file, info.Size())
}

// Apply the given handler to each entry in the given zip archive. Warning! Due
// to defects in the archive/zip library interfaces, this function copies the
// entire contents of the given reader to a temporary file and deletes that
// file before returning. Make sure that any server-side components that call
// this are configured appropriately, e.g. by allocating sufficient memory to
// their ephemeral file systems or the like.
func ForEachZipEntryFromReader(handler ZipHandler, reader io.Reader) error {
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
	return ForEachZipEntryFromFile(handler, file)
}
