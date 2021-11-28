// Package reader : Reader for bytestream buf
package reader

import (
	"encoding/binary"
	"errors"
)

// Reader : Struct to store bytestream
//          and read Position
type Reader struct {
	data   []byte
	Pos    uint16
	Length uint16
}

//var errReader = errors.New("read data out-of-boundary")

// NewReader : construct function
func NewReader(b []byte, n int) *Reader {
	return &Reader{
		data:   b,
		Pos:    0,
		Length: uint16(n),
	}
}

func (r *Reader) movePos(num uint16) {
	r.data = r.data[num:]
	r.Pos += num
}

// ReadUint8 : read 1 byte
func (r *Reader) ReadUint8() (uint8, error) {
	const Length = 1
	if r.Pos+Length > r.Length {
		return 0, errors.New("read UINT8 data out-of-boundary")
	}

	d := r.data[0]
	//binary.BigEndian does not have uint8 func
	//so directly return bytes slice
	r.movePos(Length)

	return d, nil
}

// ReadUint16 : BigEndian uint16 function
func (r *Reader) ReadUint16() (uint16, error) {
	const Length = 2
	if r.Pos+Length > r.Length {
		return 0, errors.New("read UINT16 data out-of-boundary")
	}
	d := binary.BigEndian.Uint16(r.data[0:Length])
	r.movePos(Length)

	return d, nil
}

// ReadUint32 : BigEndian uint32 function
func (r *Reader) ReadUint32() (uint32, error) {
	const Length = 4
	if r.Pos+Length > r.Length {
		return 0, errors.New("read UINT32 data out-of-boundary")
	}
	d := binary.BigEndian.Uint32(r.data[0:Length])
	r.movePos(Length)

	return d, nil
}

// ReadUint64 : BigEndian uint64 function
func (r *Reader) ReadUint64() (uint64, error) {
	const Length = 8
	if r.Pos+Length > r.Length {
		return 0, errors.New("read UINT64 data out-of-boundary")
	}
	d := binary.BigEndian.Uint64(r.data[0:Length])
	r.movePos(Length)

	return d, nil
}

// ReadN : Read N-bytes, for var-Length record type
func (r *Reader) ReadN(Length uint16) ([]byte, error) {
	if Length > 0 {
		if r.Pos+Length > r.Length {
			return []byte{}, errors.New("read N-Bytes data out-of-boundary")
		}
		d := r.data[0:Length]
		r.movePos(Length)

		return d, nil
	}

	return []byte{}, errors.New("read N-Bytes data out-of-boundary,Length Zero")

}

// FetchUint16 : BigEndian uint16 function
func (r *Reader) FetchUint16() (uint16, error) {
	const Length = 2
	if r.Pos+Length > r.Length {
		return 0, errors.New("fetch UINT16 data out-of-boundary")
	}
	d := binary.BigEndian.Uint16(r.data[0:Length])

	return d, nil
}

// FetchUint32 : BigEndian uint32 function
func (r *Reader) FetchUint32() (uint32, error) {
	const Length = 4
	if r.Pos+Length > r.Length {
		return 0, errors.New("fetch UINT32 data out-of-boundary")
	}
	d := binary.BigEndian.Uint32(r.data[0:Length])

	return d, nil
}

// FetchUint64 : BigEndian uint64 function
func (r *Reader) FetchUint64() (uint64, error) {
	const Length = 8
	if r.Pos+Length > r.Length {
		return 0, errors.New("fetch UINT64 data out-of-boundary")
	}
	d := binary.BigEndian.Uint64(r.data[0:Length])

	return d, nil
}

// FetchN : Return N bytes, but does not move Pos
func (r *Reader) FetchN(Length uint16) ([]byte, error) {
	if r.Pos+Length > r.Length {
		return []byte{}, errors.New("fetch N-Bytes data out-of-boundary")
	}
	d := r.data[0:Length]
	return d, nil
}
