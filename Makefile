NAME := ai-radar
WAYBAR_DIR := ~/.config/waybar/scripts
CLAUDE_DIR := ~/.claude/_scripts

build:
	go build -o $(WAYBAR_DIR)/$(NAME) .
	mkdir -p $(CLAUDE_DIR)
	cp cmd.sh $(CLAUDE_DIR)/$(NAME).sh
	chmod +x $(CLAUDE_DIR)/$(NAME).sh
