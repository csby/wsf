package types

import (
	"encoding/hex"
	"fmt"
	"hash/adler32"
	"strings"
)

type HttpPath interface {
	Path() string
	TokenType() int
	TokenPlace() int
	IsWebSocket() bool
	IsShortenPath() bool
	RawPath() string
	TokenUI() func() []TokenUI
	TokenCreate() func(items []TokenAuth, a Assistant) (string, ErrorCode, error)

	UseShortenPath(shortenPath bool) HttpPath
	SetWebSocket(webSocket bool) HttpPath
	SetTokenPlace(tokenPlace int) HttpPath
	SetTokenType(tokenType int) HttpPath
	SetTokenUI(tokenUI func() []TokenUI) HttpPath
	SetTokenCreate(tokenCreate func(items []TokenAuth, a Assistant) (string, ErrorCode, error)) HttpPath
}

type httpPath struct {
	rawPath       string
	path          string
	tokenType     int
	tokenPlace    int
	isWebSocket   bool
	isShortenPath bool

	tokenUI     func() []TokenUI
	tokenCreate func(items []TokenAuth, a Assistant) (string, ErrorCode, error)
}

func (s *httpPath) Path() string {
	if s.isShortenPath {
		return s.toShortenUrl(s.path)
	}

	return s.path
}

func (s *httpPath) TokenType() int {
	return s.tokenType
}

func (s *httpPath) TokenPlace() int {
	return s.tokenPlace
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

func (s *httpPath) TokenUI() func() []TokenUI {
	return s.tokenUI
}

func (s *httpPath) TokenCreate() func(items []TokenAuth, a Assistant) (string, ErrorCode, error) {
	return s.tokenCreate
}

func (s *httpPath) UseShortenPath(shortenPath bool) HttpPath {
	s.isShortenPath = shortenPath
	return s
}

func (s *httpPath) SetWebSocket(webSocket bool) HttpPath {
	s.isWebSocket = webSocket
	return s
}

func (s *httpPath) SetTokenPlace(tokenPlace int) HttpPath {
	s.tokenPlace = tokenPlace
	return s
}

func (s *httpPath) SetTokenType(tokenType int) HttpPath {
	s.tokenType = tokenType
	return s
}

func (s *httpPath) SetTokenUI(tokenUI func() []TokenUI) HttpPath {
	s.tokenUI = tokenUI
	return s
}

func (s *httpPath) SetTokenCreate(tokenCreate func(items []TokenAuth, a Assistant) (string, ErrorCode, error)) HttpPath {
	s.tokenCreate = tokenCreate
	return s
}

func (s *httpPath) toShortenUrl(url string) string {
	h := adler32.New()
	_, err := h.Write([]byte(url))
	if err != nil {
		return url
	}

	return fmt.Sprintf("/%s", hex.EncodeToString(h.Sum(nil)))
}

type Path struct {
	Prefix             string
	DefaultTokenType   int
	DefaultTokenPlace  int
	DefaultShortenUrl  bool
	DefaultTokenUI     func() []TokenUI
	DefaultTokenCreate func(items []TokenAuth, a Assistant) (string, ErrorCode, error)
}

func (s *Path) New(path string) HttpPath {
	sb := &strings.Builder{}
	sb.WriteString(s.Prefix)
	sb.WriteString(path)
	rawPath := sb.String()

	hp := &httpPath{
		rawPath:       rawPath,
		path:          rawPath,
		tokenType:     s.DefaultTokenType,
		tokenPlace:    s.DefaultTokenPlace,
		isWebSocket:   false,
		isShortenPath: s.DefaultShortenUrl,
		tokenUI:       s.DefaultTokenUI,
		tokenCreate:   s.DefaultTokenCreate,
	}

	return hp
}
