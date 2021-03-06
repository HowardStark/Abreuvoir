package entry

import (
	"bytes"
	"io"

	"github.com/HowardStark/abreuvoir/util"
)

// String Entry
type String struct {
	Base
	trueValue    string
	isPersistant bool
}

// StringFromReader builds a string entry using the provided parameters
func StringFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*String, error) {
	valLen, sizeData := util.PeekULeb128(reader)
	valData := make([]byte, valLen)
	_, err := io.ReadFull(reader, valData[:])
	if err != nil {
		return nil, err
	}
	val := string(valData[:])
	persistant := (persist == flagPersist)
	value := append(sizeData, valData[:]...)
	return &String{
		trueValue:    val,
		isPersistant: persistant,
		Base: Base{
			eName:  name,
			eType:  typeString,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}, nil
}

// StringFromItems builds a string entry using the provided parameters
func StringFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *String {
	valLen, sizeLen := util.ReadULeb128(bytes.NewReader(value))
	val := string(value[sizeLen : valLen-1])
	persistant := (persist == flagPersist)
	return &String{
		trueValue:    val,
		isPersistant: persistant,
		Base: Base{
			eName:  name,
			eType:  typeString,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the value of the String
func (stringEntry *String) GetValue() interface{} {
	return stringEntry.trueValue
}

// IsPersistant returns whether or not the entry should persist beyond restarts.
func (stringEntry *String) IsPersistant() bool {
	return stringEntry.isPersistant
}

// Clone returns an identical entry
func (stringEntry *String) Clone() *String {
	return &String{
		trueValue:    stringEntry.trueValue,
		isPersistant: stringEntry.isPersistant,
		Base:         stringEntry.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the String entry
func (stringEntry *String) CompressToBytes() []byte {
	return stringEntry.Base.compressToBytes()
}
