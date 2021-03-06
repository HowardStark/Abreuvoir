package entry

import (
	"io"

	"github.com/HowardStark/abreuvoir/util"
)

// Double Entry
type Double struct {
	Base
	trueValue    float64
	isPersistant bool
}

// DoubleFromReader builds a double entry using the provided parameters
func DoubleFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*Double, error) {
	var value [8]byte
	_, err := io.ReadFull(reader, value[:])
	if err != nil {
		return nil, err
	}
	return DoubleFromItems(name, id, sequence, persist, value[:]), nil
}

// DoubleFromItems builds a double entry using the provided parameters
func DoubleFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *Double {
	val := util.BytesToFloat64(value[:8])
	persistant := (persist == flagPersist)
	return &Double{
		trueValue:    val,
		isPersistant: persistant,
		Base: Base{
			eName:  name,
			eType:  typeDouble,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the value of the Double
func (double *Double) GetValue() interface{} {
	return double.trueValue
}

// IsPersistant returns whether or not the entry should persist beyond restarts.
func (double *Double) IsPersistant() bool {
	return double.isPersistant
}

// Clone returns an identical entry
func (double *Double) Clone() *Double {
	return &Double{
		trueValue:    double.trueValue,
		isPersistant: double.isPersistant,
		Base:         double.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the Double entry
func (double *Double) CompressToBytes() []byte {
	return double.Base.compressToBytes()
}
