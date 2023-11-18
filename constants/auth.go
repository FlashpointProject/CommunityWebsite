package constants

import "github.com/FlashpointProject/CommunityWebsite/types"

const (
	RoleAdministrator = "441043545735036929"
)

func AdminRoles() []string {
	return []string{
		RoleAdministrator,
	}
}

func StaffRoles() []string {
	return []string{
		RoleAdministrator,
	}
}

func IsModerator(roles []*types.DiscordRole) bool {
	for _, r := range roles {
		if r.ID == RoleAdministrator {
			return true
		}
	}
	return false
}
