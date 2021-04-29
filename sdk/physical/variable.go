package physical

import (
	"fmt"

	"golang.org/x/crypto/blake2b"
)

type Variable struct {
	Name       string
	Implements byte
}

func (v *Variable) NewEntry(in string) *Entry {
	valbytes := []byte(in)
	return &Entry{
		Key:       v,
		Value:     valbytes,
		KeyHash:   nil,
		ValueHash: nil,
	}
}

func (v *Variable) NewHashedEntry(in string) (*Entry, error) {
	keybytes := []byte(v.Name)
	keybytes = append(keybytes, v.Implements)
	valbytes := []byte(in)
	if keyhash, err := blake2b.New(blake2b.BlockSize, keybytes); err != nil {
		if valhash, err := blake2b.New(blake2b.BlockSize, valbytes); err != nil {
			return &Entry{
				Key:       v,
				Value:     valbytes,
				KeyHash:   keyhash,
				ValueHash: valhash,
			}, nil
		}
	} else {
		return nil, err
	}
	return nil, fmt.Errorf("could not created new entry")
}
