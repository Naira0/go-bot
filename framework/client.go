package bot

import (
	"database/sql"

	dg "github.com/bwmarrin/discordgo"
)

type Client struct {
	OwnerID       string
	Session       *dg.Session
	Commands      map[string]*Command
	CommandGroups []*Group
	DefaultPrefix string
	// the function to check if a command should be allowed to run.
	Check      func(s *dg.Session, m *dg.Message, cmd *Command) bool
	On_mention bool
	Conn       *sql.DB
	Settings   *Settings
}

func (c *Client) New(token string) error {
	session, err := dg.New("Bot " + token)

	if err != nil {
		return err
	}

	session.Identify.Intents = dg.IntentsAll

	c.Session = session

	return nil
}

func (c *Client) Get_prefix(guildID string) string {

	value, err := c.Settings.Get(guildID, "prefix")

	if value == nil && err != nil {
		return c.DefaultPrefix
	}

	prefix, _ := value["value"]

	return prefix.(string)
}

// Add_cmd will add the provided command struct to the clients command map.
// it will use the name field for the map key
func (c *Client) Add_cmd(cmd *Command) {

	if c.Commands == nil {
		c.Commands = make(map[string]*Command)
	}

	if cmd.Name == "" {
		panic("You must provide a Name field when adding a command to the client")
	}

	c.Commands[cmd.Name] = cmd

	if cmd.Group != nil {
		cmd.Group.Commands = append(cmd.Group.Commands, cmd)
	}
}

func (c *Client) Add_group(group *Group) *Group {
	c.CommandGroups = append(c.CommandGroups, group)
	return group
}
