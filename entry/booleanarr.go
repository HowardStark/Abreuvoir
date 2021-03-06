package entry

import "io"

// BooleanArr Entry
type BooleanArr struct {
	Base
	trueValue    []bool
	isPersistant bool
}

// BooleanArrFromReader builds a BooleanArr entry using the provided parameters
func BooleanArrFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*BooleanArr, error) {
	var tempValSize [1]byte
	_, sizeErr := io.ReadFull(reader, tempValSize[:])
	if sizeErr != nil {
		return nil, sizeErr
	}
	valSize := int(tempValSize[0])
	value := make([]byte, valSize)
	_, valErr := io.ReadFull(reader, value[:])
	if valErr != nil {
		return nil, valErr
	}
	return BooleanArrFromItems(name, id, sequence, persist, value), nil
}

// BooleanArrFromItems builds a BooleanArr entry using the provided parameters
func BooleanArrFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *BooleanArr {
	valSize := int(value[0])
	var val []bool
	for counter := 1; counter-1 < valSize; counter++ {
		tempVal := (value[counter] == boolTrue)
		val = append(val, tempVal)
	}
	persistant := (persist == flagPersist)
	return &BooleanArr{
		trueValue:    val,
		isPersistant: persistant,
		Base: Base{
			eName:  name,
			eType:  typeBooleanArr,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the trueValue
func (booleanArr *BooleanArr) GetValue() interface{} {
	return booleanArr.trueValue
}

// GetValueAtIndex returns the value at the specified index
func (booleanArr *BooleanArr) GetValueAtIndex(index int) bool {
	return booleanArr.trueValue[index]
}

// IsPersistant returns whether or not the entry should persist beyond restarts.
func (booleanArr *BooleanArr) IsPersistant() bool {
	return booleanArr.isPersistant
}

// Clone returns an identical entry
func (booleanArr *BooleanArr) Clone() *BooleanArr {
	return &BooleanArr{
		trueValue:    booleanArr.trueValue,
		isPersistant: booleanArr.isPersistant,
		Base:         booleanArr.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the BooleanArr entry
func (booleanArr *BooleanArr) CompressToBytes() []byte {
	return booleanArr.Base.compressToBytes()
}
