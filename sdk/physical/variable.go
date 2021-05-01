package physical

import (
	"github.com/momentag/ebliss/sdk/helpers/parseutils"
	"golang.org/x/crypto/blake2b"
)

type Variable struct {
	Name       string
	Implements byte
}

func (v *Variable) NewEntry(in interface{}) (Entry, bool) {
	if val, err := parseutils.ParseBytes(in); err != nil {
		return Entry{
			Key:       v,
			Value:     val,
			KeyHash:   nil,
			ValueHash: nil,
		}, true
	}
	return Entry{
		Key: v, Value: nil, KeyHash: nil, ValueHash: nil,
	}, false
}

func (v *Variable) NewHashedEntry(in string) (Entry, bool) {
	blake, err := blake2b.New512(nil)
	if err != nil {
		return Entry{}, false
	}
	keybytes := []byte(v.Name)
	keybytes = append(keybytes, v.Implements)
	valbytes := []byte(in)
	var keyhash, valhash []byte
	keyhash = blake.Sum(keybytes)
	valhash = blake.Sum(valbytes)
	return Entry{
		Key:       v,
		Value:     valbytes,
		KeyHash:   keyhash,
		ValueHash: valhash,
	}, true
}
