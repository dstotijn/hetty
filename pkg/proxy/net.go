package proxy

import (
	"errors"
	"net"
)

var ErrAlreadyAccepted = errors.New("listener already accepted")

// OnceListener implements net.Listener.
//
// Accepts a connection once and returns an error on subsequent
// attempts.
type OnceAcceptListener struct {
	c net.Conn
}

func (l *OnceAcceptListener) Accept() (net.Conn, error) {
	if l.c == nil {
		return nil, ErrAlreadyAccepted
	}

	c := l.c
	l.c = nil

	return c, nil
}

func (l *OnceAcceptListener) Close() error {
	return nil
}

func (l *OnceAcceptListener) Addr() net.Addr {
	return l.c.LocalAddr()
}

// ConnNotify embeds net.Conn and adds a channel field for notifying
// that the connection was closed.
type ConnNotify struct {
	net.Conn
	closed chan struct{}
}

func (c *ConnNotify) Close() error {
	err := c.Conn.Close()
	c.closed <- struct{}{}
	return err
}
