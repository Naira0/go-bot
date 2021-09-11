package bot

import (
	"fmt"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func Resolve_roleNames(r []string, guildID string, s *dg.Session) []string {
	var roles []string
	for _, id := range r {
		role, err := s.State.Role(guildID, id)

		if err != nil {
			fmt.Println("Error", err.Error())
			continue
		}

		roles = append(roles, role.Name)
	}

	return roles
}

// TODO make it so that it cleans roles and channels as well
func Clean(s string) string {
	cleaned := strings.Replace(s, "@everyone", "@\u200Beveryone", -1)
	cleaned = strings.Replace(s, "@here", "@\u200Bhere", -1)
	return cleaned
}

const (
	GREEN = 65280
)

func Has_perms(desired int64, has int64) bool {
	return has&desired == desired
}
