package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/theus-ortiz/pools-bot/internal/bot/handlers"
)

func MessageRouter(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "!") {
		return
	}

	cmdParts := strings.Fields(strings.ToLower(strings.TrimSpace(m.Content[1:])))
	cmd := cmdParts[0]

	if handler, ok := handlers.Commands[cmd]; ok {
		handler(s, m)
	} else {
		s.ChannelMessageSend(m.ChannelID, "❌ Comando não reconhecido. Digite `!help` para ver os comandos disponíveis.")
	}
}