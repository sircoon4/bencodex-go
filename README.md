Bencodex codec for Go
=======================

This library implements [Bencodex] serialization format which extends [Bencoding].

[Bencodex]: https://github.com/planetarium/bencodex
[Bencoding]: http://www.bittorrent.org/beps/bep_0003.html#bencoding


Usage
-----

It currently provides only the most basic encoder and decoder.  See also these methods:

 -  `bencodex.Encode(val any) ([]byte, error)`
 -  `bencodex.EncodeTo(w io.Writer, val any) error`
 -  `bencodex.Decode(b []byte) (any, error)`
 -  `bencodex.DecodeFrom(r io.Reader) (any, error)`

The result of decoding the encoded dictionary is returned to the internally defined type, `*bencodextype.Dictionary`.

But you can use the `map[string]any` type for `bencodex.Encode` as well as the `*bencodextype.Dictionary` type.

- `type Dictionary`
  - `func (d *Dictionary) Set(key any, value any)`
  - `func (d *Dictionary) Get(key any) any`
  - `func (d *Dictionary) Delete(key any)`
  - `func (d *Dictionary) Contains(key any) bool`
  - `func (d *Dictionary) Keys() []any`
  - `func (d *Dictionary) Values() []any`
  - `func (d *Dictionary) Length() int`
  - `func (d *Dictionary) CanConvertToMap() bool`
  - `func (d *Dictionary) ConvertToMap() map[string]any`
  - `func NewDictionary() *Dictionary`
  - `func NewDictionaryFromMap(m map[string]any) *Dictionary`


Example
-------

See [Example]