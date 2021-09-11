package main

import (
	bot "bot/framework"
	"errors"
	"strconv"
	"strings"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

func ping_cmd(ctx *bot.Context) error {

	msg, err := ctx.Send("pinging...")

	if err != nil {
		return err
	}

	started := time.Now()

	ctx.Edit(msg.ID, "pinging...")

	pingTime := time.Since(started)

	ctx.Edit(msg.ID, "Message: %s\nHeartbeat: %s", pingTime.String(), ctx.Client.Session.HeartbeatLatency().String())

	return nil
}

func echo_cmd(ctx *bot.Context) error {
	contents := bot.Clean(strings.Join(ctx.Args, " "))
	ctx.Send(contents)
	return nil
}

func help_cmd(ctx *bot.Context) error {

	if len(ctx.Args) >= 1 {
		group, exists := ctx.Find_group(ctx.Args[0])

		var command_fields []*dg.MessageEmbedField

		if exists {
			for _, command := range group.Commands {
				command_fields = append(command_fields, &dg.MessageEmbedField{
					Name:  command.Name,
					Value: command.Description,
				})
			}

			_, err := ctx.Send_embed(&dg.MessageEmbed{
				Title:  group.Name,
				Fields: command_fields,
				Color:  bot.GREEN,
				Footer: &dg.MessageEmbedFooter{Text: "use the cmd command to get more info about a command"}})

			return err
		}
	}

	var groups []*dg.MessageEmbedField

	for _, group := range client.CommandGroups {
		groups = append(groups, &dg.MessageEmbedField{Name: group.Name, Value: group.Description})
	}

	_, err := ctx.Send_embed(&dg.MessageEmbed{
		Title:       "Command groups",
		Description: "Use the command with the group name as an argument to see the commands in that group",
		Fields:      groups,
		Color:       bot.GREEN})

	return err
}

func cmd_cmd(ctx *bot.Context) error {

	cmd, exists := bot.Verify_cmd(ctx.Args[0], client.Commands)

	if !exists {
		return errors.New("Invalid command name provided")
	}

	fields := []*dg.MessageEmbedField{
		{Name: "Cooldown", Value: strconv.Itoa(cmd.Cooldown) + " Seconds", Inline: true},
		{Name: "Channel", Value: ctx.Channel_type(cmd.Channel), Inline: true},
	}

	if len(cmd.Alias) != 0 {
		fields = append(fields, &dg.MessageEmbedField{Name: "Aliases", Value: strings.Join(cmd.Alias, ", "), Inline: true})
	}

	if cmd.Needs_args {
		fields = append(fields, &dg.MessageEmbedField{Name: "Minimum Args", Value: strconv.Itoa(cmd.Arg_count), Inline: true})
	}

	_, err := ctx.Send_embed(&dg.MessageEmbed{
		Title:       cmd.Name,
		Description: cmd.Description,
		Fields:      fields,
	})

	return err
}

func roleNames_cmd(ctx *bot.Context) error {
	role := ctx.Find_role(ctx.Args[0])

	if role == nil {
		return errors.New("Could not resolve role")
	}

	ctx.Send("Role name: %s", role.Name)

	return nil
}

func test_cmd(ctx *bot.Context) error {
	member := ctx.Find_member(ctx.Args[0])

	if member == nil {
		return errors.New("Could not find member")
	}

	ctx.Send("Member: %s", member.User.Username)

	return nil
}

func perms_cmd(ctx *bot.Context) error {

	cmd, exists := ctx.Client.Commands[ctx.Args[0]]

	if !exists || cmd.Owner_only {
		return errors.New("invalid command name provided")
	}

	ctx.Send("Perms set for command `%s`", cmd.Name)

	return nil
}

func prefix_cmd(ctx *bot.Context) error {

	guild := ctx.Message.GuildID
	key := "prefix"

	if len(ctx.Args) >= 1 {
		if strings.ToLower(ctx.Args[0]) == "reset" {
			err := client.Settings.Delete(guild, key)

			if err != nil {
				return err
			}

			ctx.Send("I have reset my prefix")
			return nil
		}

		err := client.Settings.Set(guild, key, bot.Json{"value": ctx.Args[0]})

		if err != nil {
			return errors.New("Could not set new prefix " + err.Error())
		}

		ctx.Send("my new prefix is `%s`", ctx.Args[0])
	} else {
		prefix := ctx.Client.Get_prefix(guild)
		ctx.Send("My prefix `%s`", prefix)
	}

	return nil
}

func perm_test_cmd(ctx *bot.Context) error {
	perms, _ := client.Session.State.MessagePermissions(ctx.Message)
	ctx.Send("Perms %d", perms&dg.PermissionAdministrator == dg.PermissionAdministrator)
	return nil
}
