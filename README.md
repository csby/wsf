# wsf
Web Server Framework

**To get client real ip behind tcp proxy, 
need modify the following source code:**
## crypto/tls/conn.go
```go
type Conn struct {
	// constant
	conn     net.Conn
	isClient bool
	...
	// store the remote address behind proxy
	ProxyRemoteAddr string
}

func (c *Conn) RawConn() net.Conn {

	return c.conn
}
...
```
## net/http/server.go
```go
type Server struct {
	Addr    string  // TCP address to listen on, ":http" if empty
	Handler Handler // handler to invoke, http.DefaultServeMux if nil
	...
	// get the remote address behind proxy
	ProxyRemoteAddr func(net.Conn) string
}

...

func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	remoteAddr := ""
	if c.server.ProxyRemoteAddr != nil {
		remoteAddr = c.server.ProxyRemoteAddr(c.rwc)
		c.remoteAddr = remoteAddr;
	}
	...
	if tlsConn, ok := c.rwc.(*tls.Conn); ok {
		tlsConn.ProxyRemoteAddr = remoteAddr
    	...
        }
}

...

func (h initNPNRequest) ServeHTTP(rw ResponseWriter, req *Request) {
	...
	if h.c.ProxyRemoteAddr != "" {
		req.RemoteAddr = h.c.ProxyRemoteAddr
	}
	h.h.ServeHTTP(rw, req)
}
```
