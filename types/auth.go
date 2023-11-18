package types

import "time"

type AuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type SessionInfo struct {
	ID        int64     `json:"id"`
	UID       string    `json:"uid"`
	ExpiresAt time.Time `json:"expires_at"`
	IpAddr    string    `json:"ip_addr"`
}

type FPFSSProfile struct {
	ID        string         `json:"id"`
	Username  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Roles     []*DiscordRole `json:"roles"`
	Color     string         `json:"color"`
}

type UserProfile struct {
	UserID    string    `json:"uid"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	Roles     []string  `json:"roles"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cookies struct {
	Login     string
	UserID    string
	Username  string
	AvatarURL string
	Roles     string
}

type ClientCredentialsRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
}

type AuthorizationCodeGrant struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

type DiscordRole struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type FlashpointDiscordUser struct {
	ID    string         `json:"id"`
	Roles []*DiscordRole `json:"roles"`
	Color string         `json:"color"`
}
