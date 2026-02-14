package sender

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type Protocol string

const (
	UDP Protocol = "udp"
	TCP Protocol = "tcp"
	TLS Protocol = "tls"
)

type Sender struct {
	conn     net.Conn
	protocol Protocol
	addr     string
}

func NewSender(proto Protocol, host string, port int) (*Sender, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	s := &Sender{
		protocol: proto,
		addr:     addr,
	}

	if err := s.connect(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Sender) connect() error {
	var err error
	timeout := 5 * time.Second

	switch s.protocol {
	case TLS:
		conf := &tls.Config{InsecureSkipVerify: true} // For testing/simulation
		dialer := &net.Dialer{Timeout: timeout}
		s.conn, err = tls.DialWithDialer(dialer, "tcp", s.addr, conf)
	case TCP:
		s.conn, err = net.DialTimeout("tcp", s.addr, timeout)
	case UDP:
		s.conn, err = net.DialTimeout("udp", s.addr, timeout)
	default:
		return fmt.Errorf("unsupported protocol: %s", s.protocol)
	}

	return err
}

func (s *Sender) Send(msg string) error {
	if s.conn == nil {
		if err := s.connect(); err != nil {
			return err
		}
	}

	// Add newline if missing, as some receivers expect it
	if len(msg) > 0 && msg[len(msg)-1] != '\n' {
		msg += "\n"
	}

	_, err := fmt.Fprintf(s.conn, "%s", msg)
	return err
}

func (s *Sender) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
