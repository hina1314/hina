package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/hina1314/hina/server/api"
	"log"
	"net"
	"strings"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrInvalidCommand = errors.New("invalid command")
	ErrEmptyData      = errors.New("empty data")
	ErrSaving         = errors.New("error saving data")
)

type Server struct {
	Addr string
	Api  api.API
	ctx  context.Context
}

func NewServer(ctx context.Context, addr string) *Server {
	return &Server{Addr: addr, Api: api.NewAPI(), ctx: ctx}
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	defer listen.Close()
	log.Printf("Listening on %s", s.Addr)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go s.serve(conn)
	}
}

func (s *Server) serve(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)

	for {
		input, err := br.ReadString('\n')
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		val, err := s.handle(input)
		if err != nil {
			log.Printf("Handle error: %v", err)
			if _, werr := conn.Write([]byte(fmt.Sprintf("ERROR: %v\n", err))); werr != nil {
				log.Printf("Write error: %v", werr)
				return
			}
			continue
		}

		if _, err := conn.Write([]byte(val + "\n")); err != nil {
			log.Printf("Write error: %v", err)
			return
		}
	}
}

func (s *Server) handle(input string) (string, error) {
	input = strings.TrimSpace(input)
	seg := strings.Fields(input)
	if len(seg) == 0 {
		return "", ErrInvalidInput
	}

	cmd := strings.ToLower(seg[0])
	args := seg[1:]

	switch cmd {
	case "set":
		return s.handleSet(args)
	case "get":
		return s.handleGet(args)
	case "del":
		return s.handleDel(args)
	case "hset":
		return s.handleHSet(args)
	case "hget":
		return s.handleHGet(args)
	case "hgetall":
		return s.handleHGetAll(args)
	case "hdel":
		return s.handleHDel(args)
	default:
		return "", ErrInvalidCommand
	}
}

func (s *Server) handleSet(args []string) (string, error) {
	if len(args) != 2 {
		return "", ErrInvalidInput
	}
	if ok := s.Api.Set(args[0], args[1]); !ok {
		return "", ErrSaving
	}
	return "OK", nil
}

func (s *Server) handleGet(args []string) (string, error) {
	if len(args) != 1 {
		return "", ErrInvalidInput
	}
	val, ok := s.Api.Get(args[0])
	if !ok {
		return "", ErrEmptyData
	}
	return val, nil
}

func (s *Server) handleDel(args []string) (string, error) {
	if len(args) != 1 {
		return "", ErrInvalidInput
	}
	s.Api.Del(args[0])
	return "OK", nil
}

func (s *Server) handleHSet(args []string) (string, error) {
	if len(args) < 3 || len(args)%2 != 1 {
		return "", ErrInvalidInput
	}
	ok := s.Api.HSet(args[0], args[1:]...)
	if !ok {
		return "", ErrSaving
	}
	return "OK", nil
}

func (s *Server) handleHGet(args []string) (string, error) {
	if len(args) != 2 {
		return "", ErrInvalidInput
	}
	val, ok := s.Api.HGet(args[1], args[2])
	if !ok {
		return "", ErrEmptyData
	}
	return val, nil
}

func (s *Server) handleHGetAll(args []string) (string, error) {
	if len(args) != 1 {
		return "", ErrInvalidInput
	}
	val, ok := s.Api.HGetAll(args[0])
	if !ok {
		return "", ErrEmptyData
	}
	return val, nil
}

func (s *Server) handleHDel(args []string) (string, error) {
	if len(args) < 1 {
		return "", ErrInvalidInput
	}
	if ok := s.Api.HDel(args[0], args[1:]...); !ok {
		return "", ErrSaving
	}
	return "OK", nil
}
