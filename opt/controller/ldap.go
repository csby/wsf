package controller

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"strings"
)

type Ldap struct {
	Enable bool
	Host   string
	Port   int
	Base   string
}

func (s *Ldap) Authenticate(account, password string) error {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	defer l.Close()

	loginName, _ := s.getUserName(account)
	err = l.Bind(loginName, password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Ldap) getUserName(account string) (loginName, samAccountName string) {
	loginName = account
	samAccountName = account

	if index := strings.LastIndex(account, "\\"); index != -1 {
		samAccountName = account[index+1:]
	} else if index := strings.Index(account, "@"); index != -1 {
		samAccountName = account[:index]
	} else {
		domain := s.getDomain()
		if domain != "" {
			loginName = fmt.Sprintf("%s@%s", account, domain)
		}
	}

	return
}

func (s *Ldap) getDomain() string {
	if s.Base == "" {
		return ""
	}

	items := strings.Split(s.Base, ",")
	itemCount := len(items)
	if itemCount < 1 {
		return ""
	}
	item := strings.Split(items[0], "=")
	if len(item) < 2 {
		return ""
	}
	sb := &strings.Builder{}
	sb.WriteString(strings.TrimSpace(item[1]))

	for index := 1; index < itemCount; index++ {
		item := strings.Split(items[index], "=")
		if len(item) < 2 {
			break
		}
		sb.WriteString(".")
		sb.WriteString(strings.TrimSpace(item[1]))
	}

	return sb.String()
}
