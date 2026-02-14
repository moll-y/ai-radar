#!/usr/bin/env bash

# tool_name defaults to notification_type when missing.
jq -c '.tool_name //= .notification_type | {hook_event_name, tool_name, notification_type}' | timeout 1 socat -u STDIN UNIX-CONNECT:/tmp/ai-radar.sock
