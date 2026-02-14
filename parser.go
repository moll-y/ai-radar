package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

type (
	Parser struct {
		IIdx int
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

	p.save(e)
	return p.write(), nil
}

// save stores an Element in the parser list.
func (p *Parser) save(e Element) {
	for i := 0; i < len(p.List); i++ {
		// If an element with the same ID already exists in p.List, it
		// is replaced in-place.
		if p.List[i].ID == e.ID {
			p.List[i] = e
			return
		}
	}

	// If no matching ID is found, the element is written at the current
	// insertion index. The index is then advanced using circular
	// semantics, so new elements overwrite older ones once the list
	// capacity is reached.
	p.List[p.IIdx] = e
	p.IIdx = (p.IIdx + 1) % len(p.List)
}

func (p *Parser) write() string {
	var b strings.Builder
	b.Grow(len(p.List) * 22)
	for i, e := range p.List {
		if e.ID == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(":")
		b.WriteString(e.ID[:2])
		b.WriteString(" = ")
		b.WriteString(e.Status)
		b.WriteString(e.Tool)
	}
	return b.String()
}
