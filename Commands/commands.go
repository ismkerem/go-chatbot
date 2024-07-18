package commands

import(
		"github.com/bwmarrin/discordgo"

)
type Command struct {
    Name        string
    Description string
    Execute     func(*discordgo.Session, *discordgo.MessageCreate, []string)
}


var commands=[]Command{
	{
        Name:        "merhaba",
        Description: "Says hello",
        Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
            s.ChannelMessageSend(m.ChannelID, "Merhaba! Nasılsın?")
        },
    },


}