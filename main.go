package main

import (
	bot "bot/framework"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	dg "github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

func init() {
	log.SetPrefix("Error: ")
	log.SetFlags(0)
}

var client bot.Client

type Config struct {
	Token    string
	Prefix   string
	Owner_id string
}

func init_tables(db *sql.DB) error {

	_, err := db.Query(`
	create table if not exists settings(id TEXT, key TEXT UNIQUE, value JSON);`)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=go_bot password=yes900x20 sslmode=disable")

	if err != nil {
		log.Fatalf("Unable to open postgres connection %s", err.Error())
	}

	err = init_tables(db)

	if err != nil {
		log.Fatalf("Could not create one or more psql tables %s", err.Error())
	}

	client.Settings = &bot.Settings{Db: db}
	client.Conn = db

	var c Config
	data, err := os.ReadFile("config.json")

	if err != nil {
		log.Fatal("Could not open json file")
	}

	err = json.Unmarshal(data, &c)

	if err != nil {
		log.Fatal("Could not parse json file")
	}

	client.OwnerID = c.Owner_id
	client.DefaultPrefix = c.Prefix
	client.On_mention = true

	err = client.New(c.Token)

	client.Session.StateEnabled = true

	if err != nil {
		log.Fatal("Client could not login", err)
	}

	register_cmds()

	client.Session.AddHandler(on_ready)
	client.Session.AddHandler(on_message)
	client.Check = check

	err = client.Session.Open()

	if err != nil {
		client.Session.Close()
		log.Fatal("Error opening websocket connection ", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	client.Session.Close()
}

func check(s *dg.Session, m *dg.Message, cmd *bot.Command) bool {
	return true
}

func on_message(s *dg.Session, m *dg.MessageCreate) {
	_, err := bot.Eval_cmd(&client, m.Message)

	if err != nil {
		s.ChannelMessageSendEmbed(m.Message.ChannelID, &dg.MessageEmbed{
			Title:       "Error",
			Description: "```" + err.Error() + "```",
			Color:       16737894,
		})
	}
}

func on_ready(s *dg.Session, r *dg.Ready) {
	s.UpdateGameStatus(0, "Running go command framework test")
	fmt.Printf("%s is online!\n", s.State.User.Username)
}

func register_cmds() {

	general := client.Add_group(&bot.Group{
		Name:        "general",
		Description: "General commands",
		Cooldown:    3,
	})

	settings := client.Add_group(&bot.Group{
		Name:        "settings",
		Description: "Commands for configuring the bot",
		Cooldown:    10,
	})

	info := client.Add_group(&bot.Group{
		Name:        "info",
		Description: "Information based commands",
		Cooldown:    3,
	})

	client.Add_cmd(&bot.Command{
		Name:        "ping",
		Alias:       []string{"p"},
		Description: "Shows bot response time",
		Cooldown:    10,
		User_perms:  dg.PermissionSendMessages,
		Group:       general,
		Run:         ping_cmd,
	})

	client.Add_cmd(&bot.Command{
		Name:        "echo",
		Alias:       []string{"repeat", "say"},
		Description: "Repeats back the arguments",
		Cooldown:    3,
		Needs_args:  true,
		Arg_count:   1,
		Group:       general,
		Run:         echo_cmd,
	})

	client.Add_cmd(&bot.Command{
		Name:        "help",
		Alias:       []string{"h"},
		Description: "Info about all commands",
		Group:       info,
		Cooldown:    3,
		Run:         help_cmd,
	})

	client.Add_cmd(&bot.Command{
		Name:        "cmd",
		Description: "Gets info about a specified command",
		Cooldown:    3,
		Group:       info,
		Needs_args:  true,
		Arg_count:   1,
		Run:         cmd_cmd,
	})

	client.Add_cmd(&bot.Command{
		Name:        "role",
		Alias:       []string{"r"},
		Description: "Gets the name of a role by a resolvable",
		Roles:       []string{"589777016895701002"},
		Cooldown:    12,
		Needs_args:  true,
		Arg_count:   1,
		Group:       general,
		Run:         roleNames_cmd,
	})

	client.Add_cmd(&bot.Command{
		Name:        "test",
		Alias:       []string{"t"},
		Description: "testing command",
		Owner_only:  true,
		Needs_args:  true,
		Arg_count:   1,
		Run:         test_cmd,
	})

	client.Add_cmd(&bot.Command{
		Name:        "perms",
		Description: "Allows you to change command permissions",
		Group:       settings,
		Needs_args:  true,
		Arg_count:   1,
		Cooldown:    5,
		Channel:     bot.GUILD,
		Run:         perms_cmd,
	})

	client.Add_cmd(&bot.Command{
		Name:        "prefix",
		Description: "Allows you to change the server prefix of the bot",
		User_perms:  dg.PermissionManageServer,
		Group:       settings,
		Channel:     bot.GUILD,
		Run:         prefix_cmd,
	})

	client.Add_cmd(&bot.Command{
		Name: "pt",
		Run:  perm_test_cmd,
	})
}
