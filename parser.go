package main

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	word = []byte(`"tool_name":"`)
	cmds = map[byte]string{
		'B': "Batch",
		'W': "Write",
		'E': "Edit",
		'R': "Read",
	}
)

func parse(b []byte) (string, error) {
	i := bytes.Index(b, word)
	if i < 0 || len(b)-i < len(word) {
		return "", errors.New("index out of bounds")
	}

	v, ok := cmds[b[i+len(word)]]
	if !ok {
		return "", fmt.Errorf("command not found: %q", word)
	}

	return v, nil
}
