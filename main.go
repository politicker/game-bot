package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(voiceStateUpdate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func listChannels(s *discordgo.Session, guildID string) map[string]string {
	results := make(map[string]string)

	channels, err := s.GuildChannels(guildID)
	if err != nil {
		fmt.Println("error fetching guild channels,", err)
		return nil
	}

	for _, channel := range channels {
		results[channel.ID] = channel.Name
	}

	return results
}

// This function will be called (due to AddHandler above) every time a new
// voice state is created or a voice state changes.
func voiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	// Ignore if the voice state is for a voice channel leave/join to the same channel
	if vsu.BeforeUpdate != nil && vsu.VoiceState.ChannelID == vsu.BeforeUpdate.ChannelID {
		return
	}

	channels := listChannels(s, vsu.GuildID)

	channelAssocs := map[string]string{
		"game-time": "games",
		"riff-time": "riffing",
	}

	// If the user joined a new voice channel (ChannelID is not empty)
	if vsu.VoiceState.ChannelID != "" {
		channelID := vsu.ChannelID
		channelName := channels[channelID]

		chatChannelName := channelAssocs[channelName]

		if channelAssoc == "" {
			fmt.Println("No channel association found for channel ID", channelID)
			return
		}

		user, _ := s.User(vsu.VoiceState.UserID)
		s.ChannelMessageSend(channelAssoc, fmt.Sprintf("%s has joined the voice channel.", user.Username))
	}
}
