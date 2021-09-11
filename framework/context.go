package bot

import (
	"fmt"
	"regexp"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// Context stores the relevent command information
type Context struct {
	Args    []string
	Message *dg.Message
	Client  *Client
}

// Send sends a message to the context message channel, works with format specifiers
func (c *Context) Send(content string, a ...interface{}) (*dg.Message, error) {
	return c.Client.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf(content, a...))
}

//  Edit edits the given message ids message. works with format specifiers.
func (c *Context) Edit(id, content string, a ...interface{}) (*dg.Message, error) {
	return c.Client.Session.ChannelMessageEdit(c.Message.ChannelID, id, fmt.Sprintf(content, a...))
}

func (c *Context) Send_embed(embed *dg.MessageEmbed) (*dg.Message, error) {
	return c.Client.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)
}

func (c *Context) Find_group(name string) (*Group, bool) {
	for _, group := range c.Client.CommandGroups {
		if group.Name == strings.ToLower(name) {
			return group, true
		}
	}

	return nil, false
}

// Find_role gets a role struct from a string
// string must be either the role id or the role name
func (c *Context) Find_role(data string) *dg.Role {

	guildID := c.Message.GuildID

	role, err := c.Client.Session.State.Role(guildID, data)

	if err == nil {
		return role
	}

	roles, err := c.Client.Session.GuildRoles(guildID)

	if err != nil {
		return nil
	}

	for _, r := range roles {
		if strings.ToLower(r.Name) == strings.ToLower(data) {
			return r
		}
	}

	return nil
}

// currently it only works with id cause discordgo doesnt wanna fetch the guild members
func (c *Context) Find_member(data string) *dg.Member {

	guildID := c.Message.GuildID

	id_matched, err := regexp.MatchString(`\d+`, data)

	if id_matched && err == nil {
		member, err := c.Client.Session.GuildMember(guildID, data)

		if err == nil {
			return member
		}
	}

	mention_matched, err := regexp.MatchString(`^<@!*\d+>$`, data)

	if mention_matched && err == nil {
		for _, mention := range c.Message.Mentions {

			if "<@"+mention.ID+">" == data || "<@!"+mention.ID+">" == data {
				member, err := c.Client.Session.GuildMember(guildID, mention.ID)
				if err == nil {
					return member
				}
			}
		}
	}

	guild, err := c.Client.Session.State.Guild(guildID)

	if err != nil {
		return nil
	}

	for _, m := range guild.Members {
		lower := strings.ToLower(data)
		if strings.ToLower(m.User.Username) == lower || lower == strings.ToLower(m.User.String()) {
			return m
		}
	}

	return nil
}

func (c *Context) Channel_type(channel int) string {

	var value string

	switch channel {
	case 0:
		value = "Any"
		break
	case 1:
		value = "Guild"
		break
	case 2:
		value = "DM"
		break
	}

	return value
}
