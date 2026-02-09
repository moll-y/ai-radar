package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"strings"
)

type (
	Parser struct {
		List *list.List
	}

	Event struct {
		ID   string `json:"session_id"`
		Name string `json:"hook_event_name"`
		Tool string `json:"tool_name"`
	}
)

func (p *Parser) Parse(buf []byte) (string, error) {
	var event Event
	if err := json.Unmarshal(buf, &event); err != nil {
		return "", err
	}

	return p.concat(event), nil
}

// concat formats the event history as a single, space-separated string,
// ordered from oldest to newest (left → right), for example:
//
//	abc(tool1) def(tool2) ghi(tool3)
//
// The left-most entry represents the oldest event in p.List, and the
// right-most entry represents the most recent event.
func (p *Parser) concat(event Event) string {
	var b strings.Builder

	// PreToolUse         = Runs before processing a tool call.
	// PostToolUse        = Runs when a tool has executed successfully.
	// PostToolUseFailure = Runs when a tool execution fails.
	status := map[string]string{
		"PreToolUse":         "-",
		"PostToolUse":        "+",
		"PostToolUseFailure": "!",
	}

	for e := p.List.Front(); e != nil; {
		tmp := e.Next()
		if v, ok := e.Value.(Event); !ok || v.ID == event.ID {
			// Drop non-Event entries and any prior occurrence of
			// `event`, so `event` can be re-added once at the end
			// as the most recent entry.
			p.List.Remove(e)
		} else {
			b.WriteString(fmt.Sprintf("%.3s(%s%s) ", v.ID, status[v.Name], v.Tool))
		}
		e = tmp
	}

	// Append `event` as the most recent entry.
	b.WriteString(fmt.Sprintf("%.3s(%s%s)", event.ID, status[event.Name], event.Tool))
	p.List.PushBack(event)

	return b.String()
}
