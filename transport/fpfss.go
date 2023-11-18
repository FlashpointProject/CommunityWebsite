package transport

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/FlashpointProject/CommunityWebsite/config"
	"github.com/FlashpointProject/CommunityWebsite/types"
)

type FpfssToken struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type Fpfss struct {
	token       *FpfssToken
	oauthConfig *config.OauthConfig
	apiUrl      string
}

func NewFpfss(oauthConfig *config.OauthConfig, apiUrl string) (*Fpfss, error) {
	r := &Fpfss{
		oauthConfig: oauthConfig,
		apiUrl:      apiUrl,
	}
	err := r.getNewToken()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (f *Fpfss) WithFpfss(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), CommCtxKeys.FPFSS, f))
		handler.ServeHTTP(w, r)
	})
}

func (f *Fpfss) GetToken() (string, error) {
	// Validate token expiration
	var err error
	if f.token != nil {
		if f.token.ExpiresAt.After(time.Now()) {
			// Token already expired, get a new one
			err = f.getNewToken()
			if err != nil {
				return "", err
			}
		}
		if f.token.ExpiresAt.After(time.Now().Add((-60 * time.Minute))) {
			// If expiring within the next hour, get a new one
			err = f.getNewToken()
			if err != nil {
				return "", err
			}
		}
		return f.token.AccessToken, nil
	} else {
		err = f.getNewToken()
		if err != nil {
			return "", err
		}
		return f.token.AccessToken, nil
	}
}

func (f *Fpfss) GetUserRoles(uid string) (*types.FlashpointDiscordUser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/server-user/%s", f.apiUrl, uid), nil)
	if err != nil {
		return nil, err
	}
	token, err := f.GetToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var user *types.FlashpointDiscordUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (f *Fpfss) getNewToken() error {
	f.token = nil
	// Get new token
	authStr := fmt.Sprintf("%s:%s", f.oauthConfig.FpfssClientID, f.oauthConfig.FpfssClientSecret)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authStr))
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", f.oauthConfig.FpfssClientScope)
	dataStr := data.Encode()
	req, err := http.NewRequest("POST", f.oauthConfig.FpfssTokenEndpoint, bytes.NewReader([]byte(dataStr)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedAuth))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// get response body
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to get access token: %d - %s", resp.StatusCode, resp.Status)
		} else {
			return fmt.Errorf("failed to get access token: %d - %s - %s", resp.StatusCode, resp.Status, string(msg))
		}
	}
	var tokenRes *types.AuthTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenRes)
	if err != nil {
		return err
	}
	f.token = &FpfssToken{
		AccessToken: tokenRes.AccessToken,
		TokenType:   tokenRes.TokenType,
		ExpiresAt:   time.Now().Add(time.Duration(tokenRes.ExpiresIn) * time.Second),
	}

	return nil
}

func (f *Fpfss) GetGames(ids []string) ([]*types.FpfssGame, error) {
	data, err := json.Marshal(map[string]interface{}{"game_ids": ids})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/games/fetch", f.apiUrl), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	token, err := f.GetToken()
	if err != nil {
		return nil, err
	}
	fmt.Println("token: ", token)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to get response: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	var respData *types.ResponseFpfssGamesFetch
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		fmt.Printf("failed to get response: %s", err.Error())
		return nil, err
	}
	return respData.Games, nil
}

func (f *Fpfss) GetGame(id string) (*types.FpfssGame, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/game/%s", f.apiUrl, id), nil)
	if err != nil {
		return nil, err
	}
	token, err := f.GetToken()
	if err != nil {
		return nil, err
	}
	fmt.Println("token: ", token)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get game: %d - %s", resp.StatusCode, resp.Status)
	}
	var game *types.FpfssGame
	err = json.NewDecoder(resp.Body).Decode(&game)
	if err != nil {
		return nil, fmt.Errorf("failed to decode game: %w", err)
	}
	return game, nil
}
