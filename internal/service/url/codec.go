package url

import (
	"strings"

	xerrors "github.com/domovonok/url-shortener/internal/errors"
)

const (
	codeLength   = 10
	codeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	codeBase     = int64(len(codeAlphabet))
	codeSpace    = 984930291881790849 // 63^10
)

func idToCode(id int64) (string, error) {
	if id < 0 || id >= codeSpace {
		return "", xerrors.ErrInvalidCode
	}
	var code [codeLength]byte
	for i := codeLength - 1; i >= 0; i-- {
		code[i] = codeAlphabet[id%codeBase]
		id /= codeBase
	}
	return string(code[:]), nil
}

func codeToID(code string) (int64, error) {
	if len(code) != codeLength {
		return 0, xerrors.ErrInvalidCode
	}
	var id int64
	for i := 0; i < codeLength; i++ {
		index := strings.IndexByte(codeAlphabet, code[i])
		if index == -1 {
			return 0, xerrors.ErrInvalidCode
		}
		id = id*codeBase + int64(index)
	}
	return id, nil
}
