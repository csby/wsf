package types

import (
	"container/list"
	"encoding/json"
	"sync"
)

const (
	WSOptUserLogin  = 101 // 用户登陆
	WSOptUserLogout = 102 // 用户注销

	WSOptSiteUpload      = 111 // 上传并发布后台服务管理网站
	WSDocSiteUpload      = 112 // 上传并发布后台接口文档网站
	WSRootSiteUploadFile = 113 // 根站点-上传文件
	WSRootSiteDeleteFile = 114 // 根站点-删除文件
	WSWebappSiteUpload   = 115 // 上传并发布后应用网站
	WSWebappSiteDelete   = 116 // 删除应用网站
	WSCustomSiteUpload   = 119 // 上传并发布自定义网站
)

type SocketMessage struct {
	ID   int         `json:"id" note:"消息标识"`
	Data interface{} `json:"data" note:"消息内容, 结构随id而定"`
}

func (s *SocketMessage) GetData(v interface{}) error {
	data, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

type SocketChannel interface {
	Token() *Token
	Container() SocketChannelCollection
	Write(message *SocketMessage)
	Read() <-chan *SocketMessage

	getElement() *list.Element
	close()
}

type innerSocketChannel struct {
	channel   chan *SocketMessage
	element   *list.Element
	container *innerSocketChannelCollection
	token     *Token
}

func (s *innerSocketChannel) Token() *Token {
	return s.token
}

func (s *innerSocketChannel) Container() SocketChannelCollection {
	return s.container
}

func (s *innerSocketChannel) Write(message *SocketMessage) {
	select {
	case s.channel <- message:
	default:
	}
}

func (s *innerSocketChannel) Read() <-chan *SocketMessage {
	return s.channel
}

func (s *innerSocketChannel) getElement() *list.Element {
	return s.element
}

func (s *innerSocketChannel) close() {
	close(s.channel)
}

type SocketChannelCollection interface {
	OnlineUsers() []*OnlineUser
	SetListener(newChannel, removeChannel func(channel SocketChannel))
	NewChannel(token *Token) SocketChannel
	Remove(channel SocketChannel)
	Write(message *SocketMessage, token *Token)
	AddReader(reader func(message *SocketMessage, channel SocketChannel))
	Read(message *SocketMessage, channel SocketChannel)
	AddFilter(filter func(message *SocketMessage, channel SocketChannel, token *Token) bool)
}

func NewSocketChannelCollection() SocketChannelCollection {
	instance := &innerSocketChannelCollection{}
	instance.channels = list.New()
	instance.readers = make([]func(message *SocketMessage, channel SocketChannel), 0)
	instance.filters = make([]func(message *SocketMessage, channel SocketChannel, token *Token) bool, 0)

	return instance
}

type innerSocketChannelCollection struct {
	sync.RWMutex

	channels       *list.List
	readers        []func(message *SocketMessage, channel SocketChannel)
	filters        []func(message *SocketMessage, channel SocketChannel, token *Token) bool
	newListener    func(channel SocketChannel)
	removeListener func(channel SocketChannel)
}

func (s *innerSocketChannelCollection) OnlineUsers() []*OnlineUser {
	s.Lock()
	defer s.Unlock()

	tokens := make(map[string]int)
	users := make([]*OnlineUser, 0)

	for e := s.channels.Front(); e != nil; {
		ev, ok := e.Value.(SocketChannel)
		if !ok {
			break
		}

		token := ev.Token()
		if token != nil {
			_, ok := tokens[token.ID]
			if !ok {
				tokens[token.ID] = 0
				user := &OnlineUser{}
				user.CopyFrom(token)
				user.LoginDuration = user.LoginTime.Duration()

				users = append(users, user)
			}
		}

		e = e.Next()
	}

	return users
}

func (s *innerSocketChannelCollection) SetListener(newChannel, removeChannel func(channel SocketChannel)) {
	s.newListener = newChannel
	s.removeListener = removeChannel
}

func (s *innerSocketChannelCollection) NewChannel(token *Token) SocketChannel {
	s.Lock()
	defer s.Unlock()

	instance := &innerSocketChannel{container: s}
	instance.channel = make(chan *SocketMessage, 1024)
	instance.element = s.channels.PushBack(instance)
	instance.token = token
	if token != nil {
		token.Usage++
	}

	if s.newListener != nil {
		s.newListener(instance)
	}

	return instance
}

func (s *innerSocketChannelCollection) Remove(channel SocketChannel) {
	if channel == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	token := channel.Token()
	if token != nil {
		token.Usage--
	}

	if s.removeListener != nil {
		s.removeListener(channel)
	}

	s.channels.Remove(channel.getElement())
	channel.close()
}

func (s *innerSocketChannelCollection) Write(message *SocketMessage, token *Token) {
	s.Lock()
	defer s.Unlock()

	for e := s.channels.Front(); e != nil; {
		ev, ok := e.Value.(SocketChannel)
		if !ok {
			return
		}

		if !s.filter(message, ev, token) {
			ev.Write(message)
		}

		e = e.Next()
	}
}

func (s *innerSocketChannelCollection) AddReader(reader func(message *SocketMessage, channel SocketChannel)) {
	if reader == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.readers = append(s.readers, reader)
}

func (s *innerSocketChannelCollection) Read(message *SocketMessage, channel SocketChannel) {
	count := len(s.readers)
	for i := 0; i < count; i++ {
		reader := s.readers[0]
		if reader == nil {
			continue
		}

		go func(read func(message *SocketMessage, channel SocketChannel), msg *SocketMessage, ch SocketChannel) {
			defer func() {
				if err := recover(); err != nil {
				}
			}()

			read(msg, ch)
		}(reader, message, channel)
	}
}

func (s *innerSocketChannelCollection) AddFilter(filter func(message *SocketMessage, channel SocketChannel, token *Token) bool) {
	if filter == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.filters = append(s.filters, filter)
}

func (s *innerSocketChannelCollection) filter(message *SocketMessage, channel SocketChannel, token *Token) bool {
	count := len(s.filters)
	for i := 0; i < count; i++ {
		filter := s.filters[0]
		if filter == nil {
			continue
		}

		if filter(message, channel, token) {
			return true
		}
	}

	return false
}
