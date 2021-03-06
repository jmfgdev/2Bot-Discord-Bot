package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func msgPrefix(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	guild, err := guildDetails(m.ChannelID, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was a problem changing the prefix :( Try again please~")
		errorLog.Println("prefix guild details error", err.Error())
		return
	}

	if len(msglist) < 2 {
		s.ChannelMessageSend(m.ChannelID, "No prefix given :/")
		return
	}

	if m.Author.ID != guild.OwnerID && m.Author.ID != noah {
		s.ChannelMessageSend(m.ChannelID, "Sorry, only the owner can do this")
		return
	}

	var parts []string
	var space string
	msg := "without"

	if guild, ok := sMap.Server[guild.ID]; ok && !guild.Kicked {
		parts = trimSlice(strings.Split(strings.TrimPrefix(m.Content, c.Prefix+msglist[0]), "|"))
		if guild.Prefix != "" {
			parts = trimSlice(strings.Split(strings.TrimPrefix(m.Content, guild.Prefix+msglist[0]), "|"))
		}
		if len(parts) == 2 {
			if strings.ToLower(parts[1]) == "true" {
				space = " "
				msg = "with"
			}
			guild.Prefix = parts[0] + space
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Prefix changed to %s %s a trailing space", codeSeg(guild.Prefix), msg))
			saveServers()
		}
	}
	return
}

func msgGlobalPrefix(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if m.Author.ID != noah && len(msglist) < 2 {
		return
	}

	var space string
	var msg = "without"

	parts := trimSlice(strings.Split(strings.Join(msglist[1:], " "), "|"))

	if len(parts) == 2 {
		if strings.ToLower(parts[1]) == "true" {
			space = " "
			msg = "with"
		}

		c.Prefix = parts[0] + space

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":ok_hand: | All done! Prefix changed to %s %s trailing space!", c.Prefix, msg))
		saveConfig()
	}
	return
}
