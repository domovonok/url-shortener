package codec

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
)

func EncodeIDToCode(id int64) string {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(id))
	return base64.RawURLEncoding.EncodeToString(buf[:])
}

func DecodeCodeToID(code string) (int64, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return 0, err
	}
	if len(decoded) != 8 {
		return 0, errors.New("invalid code payload length")
	}
	return int64(binary.BigEndian.Uint64(decoded)), nil
}
