package main

import (
	"encoding/json"
)

type Hook struct {
	SessionID string `json:"session_id"`
	EventName string `json:"hook_event_name"`

	// tool_name defaults to notification_type when missing.
	ToolName string `json:"tool_name"`
}

func parse(buf []byte) (string, bool, string, error) {
	var h Hook
	if err := json.Unmarshal(buf, &h); err != nil {
		return "", false, "", err
	}

	// PostToolUse indicates the tool has finished executing.
	return h.SessionID[:3], h.EventName == "PostToolUse", h.ToolName, nil
}
