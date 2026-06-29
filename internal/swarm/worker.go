package swarm

import (
	"context"
	"fmt"
	"time"
)

// Worker is a slave node: it heartbeats to the master, listens for assignments
// on its inbox, executes them via Handle, and reports results back on the bus.
type Worker struct {
	Node    Node
	Bus     *Bus
	inbox   <-chan Message
	Handle  func(ctx context.Context, m Message) (string, error)
}

func NewWorker(node Node, bus *Bus) *Worker {
	node.Role = RoleWorker
	return &Worker{Node: node, Bus: bus, inbox: bus.Subscribe(node.ID)}
}

// Run blocks until ctx is cancelled, beating its heart and handling work.
func (w *Worker) Run(ctx context.Context, beat time.Duration, masterID string) {
	if beat <= 0 {
		beat = time.Second
	}
	hb := time.NewTicker(beat)
	defer hb.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-hb.C:
			w.Bus.Publish(Message{Kind: MsgHeartbeat, From: w.Node.ID, To: masterID})
		case m := <-w.inbox:
			if m.Kind != MsgAssign {
				continue
			}
			w.Bus.Publish(Message{Kind: MsgAck, From: w.Node.ID, To: m.From, TaskID: m.TaskID})
			out, err := "", error(nil)
			if w.Handle != nil {
				out, err = w.Handle(ctx, m)
			}
			kind := MsgResult
			if err != nil {
				kind, out = MsgEscalate, fmt.Sprintf("error: %v", err)
			}
			w.Bus.Publish(Message{Kind: kind, From: w.Node.ID, To: m.From, TaskID: m.TaskID, Body: out})
		}
	}
}
