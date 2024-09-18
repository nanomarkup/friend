<h2 align="center">A Go implementation of the <a href="https://nanomarkup.github.io/">Nano Markup language</a></h2>

## FUNCTIONS
```
func Compact(dst *bytes.Buffer, src []byte) error
```
Compact appends the nano-encoded src to dst, eliminating insignificant space characters.
```
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error
```
Indent function appends to `dst` the nano-encoded source (`src`) in an indented format. The data appended to dst does not begin with the prefix nor any indentation, to make it easier to embed inside other formatted nano-encoded data.
```
func Marshal(data any) ([]byte, error)
```
Marshal returns the encoding data for the input value.
It traverses the value recursively.
```
func MarshalIndent(data any, prefix, indent string) ([]byte, error)
```
MarshalIndent is like Marshal but applies Indent to format the output.
```
func Unmarshal(data []byte, v any) error
```
Unmarshal parses the encoded data and stores the result in v. If v is nil or not a pointer, Unmarshal returns an InvalidArgumentError.
It uses the inverse of the encodings that Marshal uses, allocating maps, slices, and pointers as necessary.

## TYPES
```
type InvalidArgumentError struct {
	Context string
	Err     error
}
```
InvalidArgumentError describes an error that occurs when an invalid argument is provided.
```    
func (e *InvalidArgumentError) Error() string
```
Error returns a string representation of the InvalidArgumentError.
```    
type InvalidEntityError struct {
	Context string
	Entity  string
	Err     error
}
```
InvalidEntityError describes an error that occurs when an attempt is made with an invalid entity.
```    
func (e *InvalidEntityError) Error() string
```
Error returns a string representation of the InvalidEntityError.
