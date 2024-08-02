package api

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
)

type HinaClient struct {
	ctx    context.Context
	conn   net.Conn
	reader *bufio.Reader
}

func NewHinaClient(ctx context.Context, addr string) (*HinaClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &HinaClient{
		ctx:    ctx,
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (h *HinaClient) sendCommand(cmd string) (string, error) {
	_, err := h.conn.Write([]byte(cmd))
	if err != nil {
		return "", fmt.Errorf("write error: %v", err)
	}

	response, err := h.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read error: %v", err)
	}

	return strings.TrimSpace(response), nil
}

func (h *HinaClient) Set(key, value string) error {
	cmd := fmt.Sprintf("set %s %s\n", key, value)
	response, err := h.sendCommand(cmd)
	if err != nil {
		return err
	}
	if response != "OK" {
		return errors.New(response)
	}
	return nil
}

func (h *HinaClient) Get(key string) (string, error) {
	cmd := fmt.Sprintf("get %s\n", key)
	return h.sendCommand(cmd)
}

func (h *HinaClient) Close() error {
	return h.conn.Close()
}
