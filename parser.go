package main

import (
	"encoding/json"
	"strings"
)

type (
	Parser struct {
		List []Element
	}

	Element struct {
		ID     string `json:"session_id"`
		Name   string `json:"hook_event_name"`
		Tool   string `json:"tool_name"`
		Status string
	}
)

func (p *Parser) Parse(buf []byte) (string, error) {
	var e Element
	if err := json.Unmarshal(buf, &e); err != nil {
		return "", err
	}

	switch e.Name {
	case "PreToolUse":
		e.Status = "-"
	case "PostToolUseFailure":
		e.Status = "!"
	case "PostToolUse":
		e.Status = "+"
	}

	return p.concat(e), nil
}

func (p *Parser) concat(e Element) string {
	i := 0
	for j := 0; j < len(p.List); j++ {
		if p.List[j].ID == e.ID {
			i = j
			break
		}
	}
	for i < len(p.List)-1 {
		p.List[i] = p.List[i+1]
		i++
	}
	p.List[len(p.List)-1] = e
	return p.buildString()
}

func (p *Parser) buildString() string {
	var b strings.Builder
	b.WriteString(" ")
	for _, e := range p.List {
		if e.ID == "" {
			continue
		}
		b.WriteString("[")
		b.WriteString(e.ID[:3])
		b.WriteString(" = ")
		b.WriteString(e.Status)
		b.WriteString(e.Tool)
		b.WriteString("] ")
	}
	return b.String()
}
