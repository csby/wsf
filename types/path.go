package types

import (
	"encoding/hex"
	"fmt"
	"hash/adler32"
	"strings"
)

type HttpPath interface {
	Path() string
	TokenKind() int
	TokenType() int
	IsWebSocket() bool
	IsShortenPath() bool
	RawPath() string
}

type httpPath struct {
	rawPath       string
	path          string
	tokenKind     int
	tokenType     int
	isWebSocket   bool
	isShortenPath bool
}

func (s *httpPath) Path() string {
	return s.path
}

func (s *httpPath) TokenKind() int {
	return s.tokenKind
}

func (s *httpPath) TokenType() int {
	return s.tokenType
}

func (s *httpPath) IsWebSocket() bool {
	return s.isWebSocket
}

func (s *httpPath) IsShortenPath() bool {
	return s.isShortenPath
}

func (s *httpPath) RawPath() string {
	return s.rawPath
}

type Path struct {
	Prefix            string
	TokenKind         int
	DefaultTokenType  int
	DefaultShortenUrl bool

	OnNewPath func(path HttpPath)
}

// params[0]: IsWebSocket([0: false; 1: true])
// params[1]: ShortenUrl(replace DefaultShortenUrl[0: false; 1: true])
// params[2]: TokenType(replace DefaultTokenType)
func (s *Path) New(path string, params ...int) HttpPath {
	sb := &strings.Builder{}
	sb.WriteString(s.Prefix)
	sb.WriteString(path)
	rawPath := sb.String()

	hp := &httpPath{
		rawPath:       rawPath,
		path:          rawPath,
		tokenKind:     s.TokenKind,
		tokenType:     s.DefaultTokenType,
		isWebSocket:   false,
		isShortenPath: s.DefaultShortenUrl,
	}

	paramLen := len(params)
	if paramLen > 0 {
		if params[0] != 0 {
			hp.isWebSocket = true
		} else {
			hp.isWebSocket = true
		}
	}
	if paramLen > 1 {
		if params[1] != 0 {
			hp.isShortenPath = true
		} else {
			hp.isShortenPath = false
		}
	}
	if paramLen > 2 {
		hp.tokenType = params[2]
	}

	if hp.isShortenPath {
		hp.path = s.toShortenUrl(rawPath)
	}

	if s.OnNewPath != nil {
		s.OnNewPath(hp)
	}
	return hp
}

func (s *Path) toShortenUrl(url string) string {
	h := adler32.New()
	_, err := h.Write([]byte(url))
	if err != nil {
		return url
	}

	return fmt.Sprintf("/%s", hex.EncodeToString(h.Sum(nil)))
}
