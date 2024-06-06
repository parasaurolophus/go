# parasaurolophus/go/utilities

```
package utilities // import "parasaurolophus/go/utilities"


FUNCTIONS

func Async[T any](asyncFunction func(T) T, in chan T, out chan T, panicHandler func(recovered any))
    Invoke asyncFunction for each value received on in, sending the result on
    out. If asyncFunction panics, invoke panicHandler and then send T's zero
    value to out.

func Fetch(url string) (readCloser io.ReadCloser, err error)
    Fetch a document from the given URL.

func FilterMapKeys[K comparable, V any](m map[K]V, keepKeys ...K) map[K]V
    Return a shallow copy of m, with only the keys specified by keepKeys.

func ForCSVReader(headersHandler CSVHeadersHandler, rowHandler CSVRowHandler, reader io.Reader) (err error)
    Apply the given handlers to each row in the given CSV file's contents.
    If headersHandler is non-nil, it will be applied to the first row of the
    CSV file and what it returns will be passed as the first argument to each
    subsequent invocation of rowHandler. If headersHandler is nil, the first
    argument to each invocation of rowHandler will be nil.

func ForEachZipEntry(handler ZipHandler, readerAt io.ReaderAt, size int64) error
    Apply the given handler to each entry in the given zip file.

func ForZipFile(handler ZipHandler, file *os.File) error
    Apply the given handler to each entry in the given zip archive.

func ForZipReader(handler ZipHandler, reader io.Reader) error
    Apply the given handler to each entry in the given zip archive.

func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V
    Return a map that combines the key / value pairs from all the given ones.
    The maps are processed in the given order. If the same key appears more than
    once, the value in the result will be the last one from the parameter list.


TYPES

type CSVHeadersHandler func(headers []string) ([]string, error)
    Type of function used to process the first row of a CSV file.

type CSVRowHandler func(headers, columns []string) error
    Type of function used to process each non-header row in a CSV file.

type Money interface {
        fmt.Scanner
        fmt.Stringer
        json.Marshaler
        json.Unmarshaler
        Value() float64
}
    Wrap a float64 in a struct which specifies the number of digits to emit for
    its fractional part when converting to a string or marshaling JSON.

    This type's methods only affect how the underlying floatint-point value
    is represented in text-based formats. It does not alter the mathematical
    precision of monetary values nor perform any scaling based on the
    denominations of particular currencies. For example, the appropriate number
    of digits to use for USD, CAD, GBP, EUR etc. is 2. The appropriate number
    of digits to use for JPY is 0. 100 JPY would be represented by the float64
    value 100.0.

func New(value float64, digits int) Money
    Create a Money structure initialized to the given values.

type ZipHandler func(*zip.File) error
    Type of function used to process each entry in a zip archive.
```
