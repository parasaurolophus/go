Copyright &copy; Kirk Rader 2024

# Convenience Functions for the Standard `archive/zip` Library

Output of `go doc -all`:

```
package zip // import "parasaurolophus/go/zip"


FUNCTIONS

func ForEachZipEntry(handler ZipHandler, readerAt io.ReaderAt, size int64) (err error)
    Apply the given handler to each entry in the given zip file. Terminate the
    loop upon first error.

func ForEachZipEntryFromFile(handler ZipHandler, file *os.File) error
    Apply the given handler to each entry in the given zip archive.

func ForEachZipEntryFromReader(handler ZipHandler, reader io.Reader) error
    Apply the given handler to each entry in the given zip archive. Warning!
    Due to defects in the archive/zip library interfaces, this function copies
    the entire contents of the given reader to a temporary file and deletes that
    file before returning. Make sure that any server-side components that call
    this are configured appropriately, e.g. by allocating sufficient memory to
    their ephemeral file systems or the like.


TYPES

type ZipHandler func(*zip.File) error
    Type of function used to process each entry in a zip archive.
```
