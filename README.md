# ai-radar

ai-radar is a custom Waybar module that provides real-time visibility into
claude-code activity. 

## Usage

Add (or merge) into `~/.config/waybar/config`:
```json
{
	"modules-right": ["custom/claudia", "..."],
	"custom/claudia": {
		"exec": "~/.config/waybar/scripts/ai-radar",
		"return-type": "json",
		"exec-on-event": true,
		"format": "{text}"
	},
}
```

Add (or merge) into `~/.claude/settings.json`
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "cmd.sh"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "cmd.sh"
          }
        ]
      }
    ],
    "Notification": [
      {
        "matcher": "idle_prompt",
        "hooks": [
          {
            "type": "command",
            "command": "cmd.sh"
          }
        ]
      }
    ]
  }
}
```
