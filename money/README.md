Copyright &copy; Kirk Rader 2024

# Money

Output of `go doc -all`:

```
package money // import "parasaurolophus/go/money"


TYPES

type Money interface {
        fmt.Scanner
        fmt.Stringer
        json.Marshaler
        json.Unmarshaler
        xml.Marshaler
        xml.Unmarshaler
        xml.UnmarshalerAttr
        GetDigits() uint
        GetValue() float64
        SetDigits(uint)
        SetValue(float64)
}
    An interface for representing monetary values in text-based formats. Note
    that this type is concerned only with the representation of a float64 in
    data-exchange formats like JSON or XML. Implementation should not be assumed
    to alter the precision of the underlying value in memory nor directly
    support localization such as currency symbols, fraction separator glyphs
    etc.

func NewMoney(value float64, digits uint) Money
    Create an instance that implements Money, initialized to the given values.
    The second parameter specifies the number of digits to emit when converting
    to text based representations, e.g. 2 for currencies like USD, EUR; 0 for
    JPY; etc.
```
