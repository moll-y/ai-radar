package main

import (
	"encoding/json"
	"strings"
)

type (
	Parser struct {
		ii   int
		list []Element
	}

	Element struct {
		ID   string `json:"session_id"`
		Name string `json:"hook_event_name"`
		Tool string `json:"tool_name"`
	}
)

func (p *Parser) Parse(buf []byte) (string, error) {
	var e Element
	if err := json.Unmarshal(buf, &e); err != nil {
		return "", err
	}

	switch e.Name {
	case "SessionStart":
		p.append(e)
	case "SessionEnd":
		p.remove(e)
	default:
		p.update(e)
	}

	return p.write(), nil
}

func (p *Parser) remove(e Element) {
	i := p.find(e.ID)
	if i == -1 {
		return
	}

	// Shift everything left to fill the gap.
	copy(p.list[i:], p.list[i+1:])

	// Mark last slot as empty so the writer ignores it.
	p.list[len(p.list)-1].ID = ""
}

func (p *Parser) append(e Element) {
	if p.update(e) {
		return
	}

	// If no matching ID is found, the element is written at the current
	// insertion index. The index is then advanced using circular
	// semantics, so new elements overwrite older ones once the list
	// capacity is reached.
	p.list[p.ii] = e
	p.ii = (p.ii + 1) % len(p.list)
}

func (p *Parser) update(e Element) bool {
	if i := p.find(e.ID); i != -1 {
		p.list[i] = e
		return true
	}

	return false
}

func (p *Parser) find(id string) int {
	for i := 0; i < len(p.list); i++ {
		if p.list[i].ID == id {
			return i
		}
	}

	return -1
}

func (p *Parser) write() string {
	var b strings.Builder
	b.Grow(len(p.list) * 22)
	for _, e := range p.list {
		if e.ID == "" {
			continue
		}

		if b.Len() > 0 {
			b.WriteRune(' ')
			b.WriteRune('|')
			b.WriteRune(' ')
		}

		b.WriteString(e.ID[:2])
		b.WriteRune(' ')
		b.WriteString(e.Tool)
	}

	if b.Len() == 0 {
		b.WriteString("π")
	}

	return b.String()
}
