package entry

import (
	"bytes"
	"errors"
	"io"

	"github.com/HowardStark/abreuvoir/util"
)

const (
	typeBoolean    byte = 0x00
	typeDouble     byte = 0x01
	typeString     byte = 0x02
	typeRaw        byte = 0x03
	typeBooleanArr byte = 0x10
	typeDoubleArr  byte = 0x11
	typeStringArr  byte = 0x12
	typeRPCDef     byte = 0x20

	flagTemporary byte = 0x00
	flagPersist   byte = 0x01
	flagReserved  byte = 0xFE

	boolFalse byte = 0x00
	boolTrue  byte = 0x01
)

var (
	// idSent is the required ID for an entry that is being created/sent from the client
	idSent = [2]byte{0xFF, 0xFF}
)

// Base is the base struct for entries.
type Base struct {
	eName  string
	eType  byte
	eID    [2]byte
	eSeq   [2]byte
	eFlag  byte
	eValue []byte
}

// BuildFromReader creates an entry using the reader passed in
func BuildFromReader(reader io.Reader) (Adapter, error) {
	nameLen, _ := util.ReadULeb128(reader)
	nameData := make([]byte, nameLen)
	_, nameErr := io.ReadFull(reader, nameData[:])
	if nameErr != nil {
		return nil, nameErr
	}
	name := string(nameData[:])
	var typeData [1]byte
	_, typeErr := io.ReadFull(reader, typeData[:])
	if typeErr != nil {
		return nil, typeErr
	}
	var idData [2]byte
	_, idErr := io.ReadFull(reader, idData[:])
	if idErr != nil {
		return nil, idErr
	}
	var seqData [2]byte
	_, seqErr := io.ReadFull(reader, seqData[:])
	if seqErr != nil {
		return nil, seqErr
	}
	var flagData [1]byte
	_, flagErr := io.ReadFull(reader, flagData[:])
	if flagErr != nil {
		return nil, flagErr
	}
	_, _ = nameData, name
	switch typeData[0] {
	case typeBoolean:
		return BooleanFromReader(name, idData, seqData, flagData[0], reader)
	case typeDouble:
		return DoubleFromReader(name, idData, seqData, flagData[0], reader)
	case typeString:
		return StringFromReader(name, idData, seqData, flagData[0], reader)
	case typeRaw:
		return RawFromReader(name, idData, seqData, flagData[0], reader)
	case typeBooleanArr:
		return BooleanArrFromReader(name, idData, seqData, flagData[0], reader)
	case typeDoubleArr:
		return DoubleArrFromReader(name, idData, seqData, flagData[0], reader)
	case typeStringArr:
		return StringArrFromReader(name, idData, seqData, flagData[0], reader)
	default:
		return nil, errors.New("entry: Unknown entry type")
	}
}

// BuildFromBytes creates an entry using the data passed in.
func BuildFromBytes(data []byte) (Adapter, error) {
	nameLen, sizeLen := util.ReadULeb128(bytes.NewReader(data))
	dName := string(data[sizeLen : nameLen-1])
	dType := data[nameLen]
	dID := [2]byte{data[nameLen+1], data[nameLen+2]}
	dSeq := [2]byte{data[nameLen+3], data[nameLen+4]}
	dFlag := data[nameLen+5]
	dValue := data[nameLen+6:]
	switch dType {
	case typeBoolean:
		return BooleanFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case typeDouble:
		return DoubleFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case typeString:
		return StringFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case typeRaw:
		return RawFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case typeBooleanArr:
		return BooleanArrFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case typeDoubleArr:
		return DoubleArrFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case typeStringArr:
		return StringArrFromItems(dName, dID, dSeq, dFlag, dValue), nil
	default:
		return nil, errors.New("entry: Unknown entry type")
	}
}

func (base *Base) clone() Base {
	return *base
}

// CompressToBytes remakes the original byte slice to represent this entry
func (base *Base) compressToBytes() []byte {
	var output []byte
	nameBytes := []byte(base.eName)
	nameLen := util.EncodeULeb128(uint32(len(nameBytes)))
	output = append(output, nameLen...)
	output = append(output, nameBytes...)
	output = append(output, base.eType)
	output = append(output, base.eID[:]...)
	output = append(output, base.eSeq[:]...)
	output = append(output, base.eFlag)
	output = append(output, base.eValue...)
	return output
}
