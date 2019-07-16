package controller

import (
	"github.com/csby/wsf/server/configure"
	"github.com/csby/wsf/types"
)

type controller struct {
	types.Base

	cfg *configure.Configure

	dbToken    types.TokenDatabase
	wsChannels types.SocketChannelCollection
}

func (s *controller) createCatalog(doc types.Doc, names ...string) types.Catalog {
	root := doc.AddCatalog("管理平台接口")

	count := len(names)
	if count < 1 {
		return root
	}

	child := root
	for i := 0; i < count; i++ {
		name := names[i]
		child = child.AddChild(name)
	}

	return child
}

func (s *controller) getToken(key string) *types.Token {
	if len(key) < 1 {
		return nil
	}

	if s.dbToken == nil {
		return nil
	}

	value, ok := s.dbToken.Get(key, false)
	if !ok {
		return nil
	}

	token, ok := value.(*types.Token)
	if !ok {
		return nil
	}

	return token
}

func (s *controller) writeWebSocketMessage(token string, id int, data interface{}) bool {
	if s.wsChannels == nil {
		return false
	}

	msg := &types.SocketMessage{
		ID:   id,
		Data: data,
	}

	s.wsChannels.Write(msg, s.getToken(token))

	return true
}
