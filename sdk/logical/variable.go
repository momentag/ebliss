package logical

type Variable struct {
	Name       string
	Implements byte
}

func (v *Variable) NewEntry(in interface{}) (Entry, bool) {
	if val, err := NewValue(in, v.Implements); err == nil {
		return Entry{
			Key:   v,
			Value: &val,
		}, true
	} else {
		return Entry{}, false
	}
}
