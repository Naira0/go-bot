package bot

import (
	"errors"
	"strconv"
	"strings"
	"time"

	discord "github.com/bwmarrin/discordgo"
)

const (
	GUILD = 1
	DM    = 2
)

var cooldowns = make(map[string]time.Time)

type Group struct {
	Name        string
	Description string
	Cooldown    int
	Commands    []*Command
}

type Command struct {
	// The name of the command. the Eval_cmd function will first attempt to find the command with this field
	Name        string
	Alias       []string
	Description string
	Owner_only  bool
	User_perms  int64
	//Client_perms   []string implement later
	// The required roles for a member to be able to run commands
	Roles []string
	// if set to true a member must have all roles of the Roles field to run the command
	Roles_all bool
	// if set to true a command message must have arguments
	Needs_args bool
	// will set the minimum argument count
	Arg_count int
	// command cooldown in seconds
	Cooldown int
	// the channel type the command can run in. use either GUILD or DM constants
	Channel int
	Group   *Group
	Run     func(ctx *Context) error
}

// Evalutes the message content to see if it is a command and if it can be run.
// returns false if the command was not ran.
func Eval_cmd(client *Client, message *discord.Message) (bool, error) {

	prefix := client.Get_prefix(message.GuildID)

	if message.Author.Bot || !strings.HasPrefix(message.Content, prefix) {
		return false, nil
	}

	splited := strings.Split(message.Content[len(prefix):], " ")

	if len(splited) < 1 {
		return false, nil
	}

	cmd, cmd_exists := Verify_cmd(splited[0], client.Commands)

	if !cmd_exists {
		return false, nil
	}

	if cmd.Owner_only && client.OwnerID != message.Author.ID {
		return false, errors.New("This command can only be executed by the bot owner")
	}

	arg_len := len(splited[1:])

	if cmd.Needs_args && arg_len < cmd.Arg_count {
		return false, errors.New("Command needs a minimum of " + strconv.Itoa(cmd.Arg_count) + " Args but " + strconv.Itoa(arg_len) + " were provided")
	}

	canRun, err := can_run(cmd, client, message)

	if err != nil {
		return false, err
	}

	if canRun {

		if client.Check != nil && !client.Check(client.Session, message, cmd) {
			return false, errors.New("Pre command check failed!")
		}

		cooldown, exists := cooldowns[message.Author.ID+cmd.Name]

		if !exists {
			cooldowns[message.Author.ID+cmd.Name] = time.Now()
		} else {
			now := time.Now()
			diff := int(now.Sub(cooldown).Seconds())

			if diff <= cmd.Cooldown || (cmd.Group != nil && diff <= cmd.Group.Cooldown) {
				client.Session.ChannelMessageSend(message.ChannelID, cmd.Name+" is on cooldown")
				return false, nil
			} else {
				delete(cooldowns, message.Author.ID+cmd.Name)
			}
		}

		ctx := Context{
			Args:    splited[1:],
			Message: message,
			Client:  client,
		}

		err := cmd.Run(&ctx)

		return true, err
	}

	return false, nil
}

func Verify_cmd(cmd string, commands map[string]*Command) (*Command, bool) {

	// finds command by command name
	if _, exists := commands[cmd]; exists {
		map_val := commands[cmd]
		return map_val, true
	}

	// finds command by alias
	for _, c := range commands {
		for _, alias := range c.Alias {
			if cmd == alias {
				return c, true
			}
		}
	}

	return nil, false
}

// dumb and overly complex control flow but it works
func has_roles(desired, has []string, all bool) bool {

	found := 0

	for _, id := range desired {
		for _, has_id := range has {
			if id == has_id {
				if all {
					found++
				} else {
					return true
				}
			}
		}
	}

	if found == len(desired) {
		return true
	}

	return false
}

func can_run(cmd *Command, client *Client, message *discord.Message) (bool, error) {

	channel, err := client.Session.Channel(message.ChannelID)

	if err != nil {
		return false, err
	}

	if cmd.Channel == GUILD && channel.Type == discord.ChannelTypeDM {
		return false, errors.New("This command can only be executed inside a server")
	}
	if cmd.Channel == DM && channel.Type == discord.ChannelTypeGuildText {
		return false, errors.New("This command can only be executed inside direct messages")
	}
	if cmd.Owner_only && message.Author.ID != client.OwnerID {
		return false, errors.New("This command can only be executed by the bot owner")
	}

	perms, err := client.Session.State.MessagePermissions(message)

	if err != nil {
		return false, err
	}

	if !Has_perms(cmd.User_perms, perms) {
		return false, errors.New("You do not have the correct perms to execute this command")
	}

	if len(cmd.Roles) > 0 && !has_roles(cmd.Roles, message.Member.Roles, cmd.Roles_all) {

		roles := Resolve_roleNames(cmd.Roles, message.GuildID, client.Session)

		if len(roles) == 0 {
			return false, errors.New("You do not have the required roles to execute this command")
		}

		return false, errors.New("You must have " + strings.Join(roles, ", ") + " to use this command")
	}

	return true, nil
}
