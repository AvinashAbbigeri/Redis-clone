package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/tidwall/resp"
)

const defaultListenAddr = ":5001"

// Holds the server listening address.
type Config struct {
	ListenAddr string
}

// Encapsulates the client command and the connection.
type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]bool // Tracks connected clients.
	ln        net.Listener   // TCP listener.
	addPeerCh chan *Peer     // Channel for adding new clients.
	delPeerCh chan *Peer     // Channel for removing clients.
	quitCh    chan struct{}  // Channel for server shutdown.
	msgCh     chan Message   // Channel for processing commands.

	kv *KV // Key-Value store.
}

// Creating a new server.
// Initialize the server with default values and
// create an instance of KV.
func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		delPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

// Start the TCP listener.
// Runs the loop() as a go routine to handle events.
// acceptLoop() to listen for incoming commands.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	slog.Info("goredis server running", "listenAddr", s.ListenAddr)

	return s.acceptLoop()
}

// Handles the client commands.
// SET stores the key-value pair with optional TTL.
// GET retrieves the key-value pair, ensuring it's not expired.
// HELLO returns a basic server response.
func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case ClientCommand:
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString("OK"); err != nil {
			return err
		}
	case SetCommand:
		// Set value with optional TTL
		if err := s.kv.Set(v.key, v.val, v.ttl); err != nil {
			return err
		}
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString("OK"); err != nil {
			return err
		}
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found or expired")
		}
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString(string(val)); err != nil {
			return err
		}
	case HelloCommand:
		spec := map[string]string{
			"server": "redis",
		}
		_, err := msg.peer.Send(respWriteMap(spec))
		if err != nil {
			return fmt.Errorf("peer send error: %s", err)
		}
	}
	return nil
}

// The event loop.
// Handles the peer connections and disconnections.
// Uses select statement to handle multiple connections.
func (s *Server) loop() {
	ticker := time.NewTicker(1 * time.Second) // Cleanup expired keys every second
	defer ticker.Stop()

	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message error", "err", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			slog.Info("peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case peer := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		case <-ticker.C:
			s.kv.cleanupExpiredKeys()
		}
	}
}

// Accepts new connections and
// launches a new go routine for each connection made.
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

// Handles the client connection.
// Creates a peer instance and
// reads continuously to process client messages.
func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh, s.delPeerCh)
	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

// Reads the listenAddr flag to create a new server instance and logs if crashes.
func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddr, "listen address of the goredis server")
	flag.Parse()
	server := NewServer(Config{
		ListenAddr: *listenAddr,
	})
	log.Fatal(server.Start())
}
