package transport

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/FlashpointProject/CommunityWebsite/constants"
	"github.com/FlashpointProject/CommunityWebsite/logging"
	"github.com/FlashpointProject/CommunityWebsite/service"
	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/FlashpointProject/CommunityWebsite/utils"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

var cookies = types.Cookies{
	Login:     "login",
	UserID:    "uid",
	Username:  "username",
	AvatarURL: "avatar_url",
	Roles:     "roles",
}

type State struct {
	Nonce       string `json:"nonce"`
	RedirectURI string `json:"redirect_uri"`
}

type StateKeeper struct {
	sync.Mutex
	states            map[string]time.Time
	expirationSeconds uint64
}

// Generate generates state and returns base64-encoded form
func (sk *StateKeeper) Generate(redirectURI string) (string, error) {
	sk.Clean()
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	s := &State{
		Nonce:       u.String(),
		RedirectURI: redirectURI,
	}
	sk.Lock()
	sk.states[s.Nonce] = time.Now()
	sk.Unlock()

	j, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	b := base64.URLEncoding.EncodeToString(j)

	return b, nil
}

// Consume consumes base64-encoded state and returns destination URL
func (sk *StateKeeper) Consume(b string) (string, bool) {
	sk.Clean()
	sk.Lock()
	defer sk.Unlock()

	j, err := base64.URLEncoding.DecodeString(b)
	if err != nil {
		return "", false
	}

	s := &State{}

	err = json.Unmarshal(j, s)
	if err != nil {
		return "", false
	}

	_, ok := sk.states[s.Nonce]
	if ok {
		delete(sk.states, s.Nonce)
	}
	return s.RedirectURI, ok
}

func (sk *StateKeeper) Clean() {
	sk.Lock()
	defer sk.Unlock()
	for k, v := range sk.states {
		if v.After(v.Add(time.Duration(sk.expirationSeconds))) {
			delete(sk.states, k)
		}
	}
}

var stateKeeper = StateKeeper{
	states:            make(map[string]time.Time),
	expirationSeconds: 30,
}

func (a *App) OAuthLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, err := url.Parse(a.Conf.OauthConfig.AuthorizeEndpoint)
	if err != nil {
		utils.LogCtx(ctx).Error("failed to parse oauth authorize endpoint")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a state
	cb_dest := r.URL.Query().Get("dest")
	cb_redirect := r.URL.Query().Get("redirect_uri")
	if cb_redirect == "" && cb_dest != "" {
		cb_redirect = cb_dest
	}
	if cb_redirect == "" {
		cb_redirect = "/"
	}

	state, err := stateKeeper.Generate(cb_redirect)
	if err != nil {
		utils.LogCtx(ctx).Error("failed to generate state")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", a.Conf.OauthConfig.ClientID)
	q.Add("redirect_uri", a.Conf.OauthConfig.Callback)
	q.Add("scope", a.Conf.OauthConfig.Scope)
	q.Add("state", state)
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func (a *App) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify state

	redirectUri, ok := stateKeeper.Consume(r.FormValue("state"))
	if !ok {
		writeError(ctx, w, perr("state does not match", http.StatusBadRequest))
		return
	}
	valid := isReturnURLLocal(redirectUri, a.Conf.HostBaseUrl)
	if !valid {
		redirectUri = "/"
	}

	authStr := fmt.Sprintf("%s:%s", a.Conf.OauthConfig.ClientID, a.Conf.OauthConfig.ClientSecret)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authStr))
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", a.Conf.OauthConfig.Callback)
	// Use access token
	req, err := http.NewRequest("POST", a.Conf.OauthConfig.TokenEndpoint, bytes.NewReader([]byte(data.Encode())))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedAuth))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get access token", http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	// Read access token response
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.LogCtx(ctx).Error("failed to get access token", resp.StatusCode)
		} else {
			utils.LogCtx(ctx).Error(fmt.Sprintf("failed to get access token: %s", string(msg)), resp.StatusCode)
		}
		writeError(ctx, w, perr("failed to get access token", http.StatusInternalServerError))
		return
	}

	var res *types.AuthTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to decode access token response", http.StatusInternalServerError))
		return
	}

	// Read user profile info from Discord
	req, err = http.NewRequest("GET", a.Conf.OauthConfig.ProfileEndpoint, nil)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", res.AccessToken))
	resp, err = client.Do(req)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get user profile", http.StatusInternalServerError))
		return
	}
	var discordUser *DiscordOauthProfile
	err = json.NewDecoder(resp.Body).Decode(&discordUser)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to decode user profile", http.StatusInternalServerError))
		return
	}
	discordUid, err := strconv.ParseInt(discordUser.User.ID, 10, 64)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to parse discord user ID", http.StatusInternalServerError))
		return
	}

	// Prefer global username
	username := discordUser.User.GlobalName
	if username == "" {
		username = discordUser.User.Username
	}

	// Form avatar URL
	avatar := fmt.Sprintf("embed/avatars/%d.png", (discordUid>>22)%6) // Magic numbers discord pls
	if discordUser.User.Avatar != "" {
		avatar = fmt.Sprintf("avatars/%s/%s.png", discordUser.User.ID, discordUser.User.Avatar)
	}

	// Form into FPFSS compatible profile
	fpfssUser := &types.FPFSSProfile{
		ID:        discordUser.User.ID,
		Username:  username,
		AvatarURL: fmt.Sprintf("https://cdn.discordapp.com/%s", avatar),
		Roles:     []*types.DiscordRole{},
	}

	// Get user roles from FPFSS
	flashpointUser, err := a.Fpfss.GetUserRoles(fpfssUser.ID)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get user roles", http.StatusInternalServerError))
		return
	}
	fpfssUser.Roles = flashpointUser.Roles
	fpfssUser.Color = flashpointUser.Color
	if fpfssUser.Color == "#000000" {
		fpfssUser.Color = ""
	}

	ipAddr := logging.RequestGetRemoteAddress(r)

	// Make session, discard Discord auth token (not needed anymore)
	authToken, err := a.Service.SaveUser(ctx, fpfssUser, ipAddr)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to save user", http.StatusInternalServerError))
		return
	}

	if err := a.CC.SetSecureCookie(w, cookies.Login, service.MapAuthToken(authToken), (int)(a.Conf.SessionExpirationSeconds)); err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to set cookie", http.StatusInternalServerError))
		return
	}

	// Encode roles
	roles := make([]string, len(fpfssUser.Roles))
	for i, role := range fpfssUser.Roles {
		roles[i] = role.ID
	}
	SetCookie(w, cookies.UserID, fpfssUser.ID, (int)(a.Conf.SessionExpirationSeconds))
	SetCookie(w, cookies.Username, fpfssUser.Username, (int)(a.Conf.SessionExpirationSeconds))
	SetCookie(w, cookies.AvatarURL, fpfssUser.AvatarURL, (int)(a.Conf.SessionExpirationSeconds))
	SetCookie(w, cookies.Roles, strings.Join(roles, ","), (int)(a.Conf.SessionExpirationSeconds))

	// Redirect to final page
	http.Redirect(w, r, redirectUri, http.StatusFound)
}

func (a *App) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)

	profile, err := a.Service.GetUser(ctx, uid)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get profile", http.StatusInternalServerError))
		return
	}

	if profile == nil {
		writeError(ctx, w, perr("profile not found", http.StatusNotFound))
		return
	}

	writeResponse(ctx, w, profile, http.StatusOK)
}

func (a *App) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	uid := params[constants.ResourceKeyUserID]

	profile, err := a.Service.GetUser(ctx, uid)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get profile", http.StatusInternalServerError))
		return
	}

	if profile == nil {
		writeError(ctx, w, perr("profile not found", http.StatusNotFound))
		return
	}

	writeResponse(ctx, w, profile, http.StatusOK)
}

type DiscordOauthProfile struct {
	User *DiscordUser `json:"user"`
}

type DiscordUser struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	GlobalName string `json:"global_name"`
}
