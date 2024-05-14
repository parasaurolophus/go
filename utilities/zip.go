package utilities

import (
	"archive/zip"
	"bytes"
	"io"
)

// Return an io.ReadCloser list, one for each entry in the given zip archive.
//
// Warning: The zip.NewReader() constructor requires an io.ReaderAt which is
// incompatible with http.Response.Body() and so, as a work-around, we buffer
// the entire input in memory. Plan accordingly when provisioning to run this
// code as a service!
func Unzip(reader io.ReadCloser) ([]io.ReadCloser, error) {
	defer reader.Close()
	///////////////////////////////////////////////////////////////////////////
	// TODO: investigate better ways to do this than to read the entire input
	// file into memory (but there may not be given the design defects in the
	// standard zip library)
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	bufferReader := bytes.NewReader(buffer)
	///////////////////////////////////////////////////////////////////////////
	z, err := zip.NewReader(bufferReader, int64(len(buffer)))
	if err != nil {
		return nil, err
	}
	result := []io.ReadCloser{}
	for _, f := range z.File {
		contents, err := f.Open()
		if err != nil {
			for _, rc := range result {
				rc.Close()
			}
			return nil, err
		}
		result = append(result, contents)
	}
	return result, nil
}
