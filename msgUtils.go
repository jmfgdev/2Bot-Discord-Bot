package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var mem runtime.MemStats

/*
	These are usually short commands that dont warrant their own file
	or are only for me, the creator..usually
*/

func msgSetGame(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if m.Author.ID != noah && len(msglist) < 2 {
		return
	}

	game := strings.Join(msglist[1:], " ")

	err := s.UpdateStatus(0, game)
	if err != nil {
		log(true, "Game change error", err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, ":ok_hand: | Game changed successfully!")

	c.Game = game
	saveConfig()
	return
}

func msgInfo(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	ct1, _ := getCreationTime(s.State.User.ID)
	creationTime := ct1.Format(time.UnixDate)[:10]
	runtime.ReadMemStats(&mem)

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: 0,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    s.State.User.Username,
			IconURL: discordgo.EndpointUserAvatar(s.State.User.ID, s.State.User.Avatar),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Brought to you by 2Bot :)\nLast Bot reboot: " + lastReboot + " GMT",
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Bot Name:", Value: codeBlock(s.State.User.Username), Inline: true},
			{Name: "Creator:", Value: codeBlock("Strum355#1180"), Inline: true},
			{Name: "Creation Date:", Value: codeBlock(creationTime), Inline: true},
			{Name: "Global Prefix:", Value: codeBlock(c.Prefix), Inline: true},
			{Name: "Programming Language:", Value: codeBlock("Go"), Inline: true},
			{Name: "Library:", Value: codeBlock("Discordgo"), Inline: true},
			{Name: "Server Count:", Value: codeBlock(strconv.Itoa(len(s.State.Guilds))), Inline: true},
			{Name: "Memory Usage:", Value: codeBlock(strconv.Itoa(int(mem.Alloc/1024/1024)) + "MB"), Inline: true},
			{Name: "My Server:", Value: "https://discord.gg/9T34Y6u\nJoin here for support amongst other things!", Inline: false},
		},
	})
	return
}

func msgListUsers(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if m.Author.ID != noah && len(msglist) < 2 {
		return
	}

	if guild, ok := c.Servers[msglist[1]]; ok && !guild.Kicked {
		s.ChannelTyping(m.ChannelID)
		var out []string
		guildDetails, err := guildDetails(msglist[1], s)
		if err != nil {
			return
		}
		for _, user := range guildDetails.Members {
			out = append(out, user.User.Username)
		}

		s.ChannelMessageSend(m.ChannelID, "Users in: "+guildDetails.Name+"\n`"+strings.Join(out, "`, `")+"`")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "2Bot isn't in that server")
	return
}

func msgAnnounce(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if m.Author.ID != noah && len(msglist) < 2 {
		return
	}
	//Discord Bots, cool kidz only, social experiment, discord go
	blacklist := []string{"110373943822540800", "272873324705742848", "244133074328092673", "118456055842734083"}
	for _, guild := range s.State.Guilds {
		if !isIn(guild.ID, blacklist) {
			if _, ok := c.Servers[guild.ID]; !ok {
				log(true, "State and config mis-match")

			}
			s.ChannelMessageSend(guild.ID, strings.Join(msglist[1:], " "))
		}
	}
}

func msgGit(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(m.ChannelID, "Check me out here https://github.com/Strum355/2Bot-Discord-Bot\nGive it star to make my creators day!")
}

func msgNSFW(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	onOrOff := map[bool]string{true: "enabled", false: "disabled"}

	guildDetails, err := guildDetails(m.ChannelID, s)
	if err != nil {
		return
	}

	if m.Author.ID != guildDetails.OwnerID || m.Author.ID != noah {
		s.ChannelMessageSend(m.ChannelID, "Sorry, only the owner can do this")
		return
	}

	if guild, ok := c.Servers[guildDetails.ID]; ok && !guild.Kicked {
		guild.Nsfw = !guild.Nsfw
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("NSFW %s", onOrOff[guild.Nsfw]))
		saveConfig()
	}
	return
}

func msgJoinMessage(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	guildDetails, err := guildDetails(m.ChannelID, s)
	if err != nil {
		return
	}

	if m.Author.ID != guildDetails.OwnerID || m.Author.ID != noah {
		s.ChannelMessageSend(m.ChannelID, "Sorry, only the owner can do this")
		return
	}

	split := trimSlice(strings.Split(strings.Join(msglist[1:], " "), "|"))

	if len(split) == 0 {
		split = append(split, msglist[1])
	}

	if len(split) > 0 {
		if guild, ok := c.Servers[guildDetails.ID]; ok && !guild.Kicked {
			if split[0] != "false" && split[0] != "true" {
				s.ChannelMessageSend(m.ChannelID, "Please say either `true` or `false` for enabling or disabling join messages~")
				return
			}

			if split[0] == "false" {
				guild.JoinMessage = [3]string{split[0]}
				saveConfig()
				s.ChannelMessageSend(m.ChannelID, "Join messages disabled! ")
				return
			}

			if len(split) != 3 {
				s.ChannelMessageSend(m.ChannelID, "Not enough info given! :/\nMake sure the command only has two `|` in it.")
				return
			}
			channelStruct, err := s.State.Channel(split[2])
			if err != nil {
				log(true, "Join message channel struct or bad channel ID?", split[2], err.Error())
				s.ChannelMessageSend(m.ChannelID, "Please give me a proper channel ID :(")
				return
			}

			if split[1] == "" {
				s.ChannelMessageSend(m.ChannelID, "No message given :/")
				return
			}

			guild.JoinMessage = [3]string{split[0], split[1], split[2]}
			saveConfig()

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Join message set to:\n%s\nin %s", split[1], channelStruct.Name))
		}
	}
}

func msgReloadConfig(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if m.Author.ID != noah {
		return
	}

	if len(msglist) < 2 {
		return
	}
	var reloaded string
	switch msglist[1] {
	case "c":
		c = &config{}
		if err := loadConfig(); err != nil {
			log(true, "Error reloading config", err.Error())
			s.ChannelMessageSend(m.ChannelID, "Error reloading config")
			return
		}
		reloaded = "config"
	case "u":
		u = &users{}
		if err := loadConfig(); err != nil {
			log(true, "Error reloading config", err.Error())
			s.ChannelMessageSend(m.ChannelID, "Error reloading config")
			return
		}
		reloaded = "users"
	}

	s.ChannelMessageSend(m.ChannelID, "Reloaded "+reloaded)
}

func msgInvite(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: 0,
		Image: &discordgo.MessageEmbedImage{
			URL: happyEmoji,
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Invite me with this link!", Value: "https://discordapp.com/oauth2/authorize?client_id=301819949683572738&scope=bot&permissions=3533824", Inline: true},
		},
	})
}

func msgPrintJSON(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if len(msglist) < 3 {
		return
	}
	switch msglist[1] {
	case "u":
		if _, ok := u.User[msglist[2]]; ok {
			var out bytes.Buffer
			err := json.Indent(&out, []byte(fmt.Sprintf("%v", *u)), "", "  ")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(fmt.Sprintf("%v", *u))
			s.ChannelMessageSend(m.ChannelID, string(out.Bytes()))
		}
	}
}
