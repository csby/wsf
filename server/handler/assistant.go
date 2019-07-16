package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/csby/security/certificate"
	"io"
	"runtime"
	"time"
)

type assistant struct {
	restart    func() error
	keys       map[string]interface{}
	log        bool
	rid        uint64
	rip        string
	enterTime  time.Time
	leaveTime  time.Time
	input      []byte
	output     []byte
	rsaPrivate *certificate.RSAPrivate
}

func (s *assistant) CanUpdate() bool {
	if s.restart == nil {
		return false
	} else if runtime.GOOS == "linux" {
		return true
	} else {
		return false
	}
}

func (s *assistant) CanRestart() bool {
	if s.restart == nil {
		return false
	} else {
		return true
	}
}

func (s *assistant) Restart() error {
	if s.restart == nil {
		return fmt.Errorf("restart not supported")
	}

	return s.restart()
}

func (s *assistant) Set(key string, val interface{}) {
	s.keys[key] = val
}

func (s *assistant) Get(key string) (interface{}, bool) {
	val, ok := s.keys[key]
	if ok {
		return val, true
	} else {
		return nil, false
	}
}

func (s *assistant) Del(key string) bool {
	_, ok := s.keys[key]
	if ok {
		delete(s.keys, key)
		return true
	} else {
		return false
	}
}

func (s *assistant) SetLog(v bool) {
	s.log = v
}

func (s *assistant) GetLog() bool {
	return s.log
}

func (s *assistant) EnterTime() time.Time {
	return s.enterTime
}

func (s *assistant) LeaveTime() time.Time {
	return s.leaveTime
}

func (s *assistant) RID() uint64 {
	return s.rid
}

func (s *assistant) RIP() string {
	return s.rip
}

func (s *assistant) SetInput(v []byte) {
	s.input = v
}

func (s *assistant) GetInput() []byte {
	return s.input
}

func (s *assistant) GetOutput() []byte {
	return s.output
}

func (s *assistant) NewGuid() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func (s *assistant) RSAPublicKey() string {
	if s.rsaPrivate == nil {
		return ""
	}

	publicKey, err := s.rsaPrivate.Public()
	if err == nil {
		keyVal, err := publicKey.ToMemory()
		if err == nil {
			return string(keyVal)
		}
	}

	return ""
}

func (s *assistant) RSADecrypt(data string) (string, error) {
	if s.rsaPrivate == nil {
		return "", fmt.Errorf("rsa private key is nll")
	}

	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	decrypted, err := s.rsaPrivate.Decrypt(buf)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
