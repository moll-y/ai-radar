package main

import (
	"encoding/json"
)

type Hook struct {
	EventName string `json:"hook_event_name"`

	// tool_name defaults to notification_type when missing.
	ToolName string `json:"tool_name"`
}

func parse(buf []byte) (bool, string, error) {
	var h Hook
	if err := json.Unmarshal(buf, &h); err != nil {
		return false, "", err
	}

	// PostToolUse indicates the tool has finished executing.
	return h.EventName == "PostToolUse", h.ToolName[:2], nil
}
