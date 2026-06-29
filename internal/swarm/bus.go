package swarm

import (
	"sync"
	"time"
)

// MsgKind enumerates the internal agent-to-agent protocol verbs that ride the
// gossip bus between the master and worker agents.
type MsgKind string

const (
	MsgHeartbeat   MsgKind = "heartbeat"
	MsgAssign      MsgKind = "assign"
	MsgAck         MsgKind = "ack"
	MsgProgress    MsgKind = "progress"
	MsgHandoff     MsgKind = "handoff"
	MsgResult      MsgKind = "result"
	MsgEscalate    MsgKind = "escalate"
	MsgGossip      MsgKind = "gossip"
	MsgBroadcast   MsgKind = "broadcast"
)

// Message is the envelope passed on the internal bus. From/To are node IDs;
// To == "*" means broadcast to the whole swarm.
type Message struct {
	ID      string            `json:"id"`
	Kind    MsgKind           `json:"kind"`
	From    string            `json:"from"`
	To      string            `json:"to"`
	TaskID  string            `json:"task_id"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
	TS      time.Time         `json:"ts"`
}

type subscriber struct {
	id string
	ch chan Message
}

// Bus is an in-process pub/sub fabric: one inbox per node plus a broadcast
// fan-out. This is the seam where a real NATS/Redis transport would slot in.
type Bus struct {
	mu      sync.RWMutex
	subs    map[string]*subscriber
	history []Message
	maxHist int
}

func NewBus() *Bus {
	return &Bus{subs: make(map[string]*subscriber), maxHist: 512}
}

// Subscribe registers a node inbox and returns its receive channel.
func (b *Bus) Subscribe(nodeID string) <-chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := &subscriber{id: nodeID, ch: make(chan Message, 64)}
	b.subs[nodeID] = s
	return s.ch
}

func (b *Bus) Unsubscribe(nodeID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if s, ok := b.subs[nodeID]; ok {
		close(s.ch)
		delete(b.subs, nodeID)
	}
}

// Publish routes a message to one inbox or fans out to every subscriber.
func (b *Bus) Publish(m Message) {
	if m.TS.IsZero() {
		m.TS = time.Now()
	}
	b.mu.Lock()
	b.history = append(b.history, m)
	if len(b.history) > b.maxHist {
		b.history = b.history[len(b.history)-b.maxHist:]
	}
	targets := make([]*subscriber, 0, len(b.subs))
	for _, s := range b.subs {
		if m.To == "*" || m.To == s.id {
			targets = append(targets, s)
		}
	}
	b.mu.Unlock()

	for _, s := range targets {
		select {
		case s.ch <- m:
		default:
			// slow consumer: drop rather than block the swarm
		}
	}
}

// Tail returns the most recent n messages for the audit timeline / TUI.
func (b *Bus) Tail(n int) []Message {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if n > len(b.history) {
		n = len(b.history)
	}
	out := make([]Message, n)
	copy(out, b.history[len(b.history)-n:])
	return out
}
