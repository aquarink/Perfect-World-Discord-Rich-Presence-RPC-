package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	OpHandshake = 0
	OpFrame     = 1
	OpClose     = 2
	OpPing      = 3
	OpPong      = 4
)

type ActivityTimestamps struct {
	Start int64 `json:"start,omitempty"`
}

type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

type ActivityButton struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

type Activity struct {
	Details    string             `json:"details,omitempty"`
	State      string             `json:"state,omitempty"`
	Timestamps ActivityTimestamps `json:"timestamps,omitempty"`
	Assets     ActivityAssets     `json:"assets,omitempty"`
	Buttons    []ActivityButton   `json:"buttons,omitempty"`
}

type FrameArgs struct {
	Pid      int       `json:"pid"`
	Activity *Activity `json:"activity"`
}

type FramePayload struct {
	Cmd   string    `json:"cmd"`
	Args  FrameArgs `json:"args"`
	Nonce string    `json:"nonce"`
}

type HandshakePayload struct {
	V        string `json:"v"`
	ClientID string `json:"client_id"`
}

type DiscordClient struct {
	ClientID string
	conn     io.ReadWriteCloser
}

func NewDiscordClient(clientID string) *DiscordClient {
	return &DiscordClient{
		ClientID: clientID,
	}
}

func openDiscordPipe() (io.ReadWriteCloser, error) {
	if runtime.GOOS == "windows" {
		for i := 0; i < 10; i++ {
			pipePath := fmt.Sprintf(`\\.\pipe\discord-ipc-%d`, i)
			f, err := os.OpenFile(pipePath, os.O_RDWR, 0)
			if err == nil {
				return f, nil
			}
		}
		return nil, fmt.Errorf("could not open discord ipc pipe (is Discord running?)")
	}

	// Linux / macOS fallback for development/testing
	envKeys := []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}
	var tempDirs []string
	for _, key := range envKeys {
		if val := os.Getenv(key); val != "" {
			tempDirs = append(tempDirs, val)
		}
	}
	tempDirs = append(tempDirs, "/tmp")

	for _, dir := range tempDirs {
		for i := 0; i < 10; i++ {
			sockPath := fmt.Sprintf("%s/discord-ipc-%d", dir, i)
			conn, err := net.Dial("unix", sockPath)
			if err == nil {
				return conn, nil
			}
		}
	}

	return nil, fmt.Errorf("could not open discord unix socket (is Discord running?)")
}

func (c *DiscordClient) Connect() error {
	conn, err := openDiscordPipe()
	if err != nil {
		return err
	}
	c.conn = conn

	// Handshake
	hs := HandshakePayload{
		V:        "1",
		ClientID: c.ClientID,
	}
	data, err := json.Marshal(hs)
	if err != nil {
		c.Close()
		return err
	}

	if err := c.writePacket(OpHandshake, data); err != nil {
		c.Close()
		return err
	}

	// Read handshake reply
	_, _, err = c.readPacket()
	if err != nil {
		c.Close()
		return err
	}

	return nil
}

func (c *DiscordClient) SetActivity(activity Activity) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to Discord")
	}

	payload := FramePayload{
		Cmd: "SET_ACTIVITY",
		Args: FrameArgs{
			Pid:      os.Getpid(),
			Activity: &activity,
		},
		Nonce: strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.writePacket(OpFrame, data)
}

func (c *DiscordClient) ClearActivity() error {
	if c.conn == nil {
		return nil
	}

	payload := FramePayload{
		Cmd: "SET_ACTIVITY",
		Args: FrameArgs{
			Pid:      os.Getpid(),
			Activity: nil,
		},
		Nonce: strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.writePacket(OpFrame, data)
}

func (c *DiscordClient) Close() {
	if c.conn != nil {
		_ = c.ClearActivity()
		_ = c.conn.Close()
		c.conn = nil
	}
}

func (c *DiscordClient) writePacket(op int32, data []byte) error {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, op); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, int32(len(data))); err != nil {
		return err
	}
	if _, err := buf.Write(data); err != nil {
		return err
	}
	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *DiscordClient) readPacket() (int32, []byte, error) {
	header := make([]byte, 8)
	if _, err := io.ReadFull(c.conn, header); err != nil {
		return 0, nil, err
	}

	op := int32(binary.LittleEndian.Uint32(header[0:4]))
	length := int32(binary.LittleEndian.Uint32(header[4:8]))

	body := make([]byte, length)
	if _, err := io.ReadFull(c.conn, body); err != nil {
		return 0, nil, err
	}

	return op, body, nil
}
