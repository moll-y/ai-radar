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

	Input struct {
		SessionID string `json:"session_id"`
		EventName string `json:"hook_event_name"`
		// tool_name defaults to notification_type when missing.
		ToolName string `json:"tool_name"`
	}
)

func (p *Parser) Parse(buf []byte) (string, error) {
	var in Input
	if err := json.Unmarshal(buf, &in); err != nil {
		return "", err
	}

	return p.concat(in), nil
}

func (p *Parser) concat(in Input) string {
	var b strings.Builder

	for e := p.List.Front(); e != nil; {
		tmp := e.Next()
		if v, ok := e.Value.(Input); ok && v.SessionID != in.SessionID {
			b.WriteString(fmt.Sprintf("%.3s(%s) ", v.SessionID, v.ToolName))
		} else {
			p.List.Remove(e)
		}
		e = tmp
	}

	p.List.PushBack(in)
	b.WriteString(fmt.Sprintf("%.3s(%s)", in.SessionID, in.ToolName))
	return b.String()
}
