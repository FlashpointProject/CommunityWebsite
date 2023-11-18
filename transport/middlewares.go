package transport

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/FlashpointProject/CommunityWebsite/constants"
	"github.com/FlashpointProject/CommunityWebsite/service"
	"github.com/FlashpointProject/CommunityWebsite/utils"
)

// UserAuthMux takes many authorization middlewares and accepts if any of them does not return error
func (a *App) UserAuthMux(next func(http.ResponseWriter, *http.Request), authorizers ...func(*http.Request, string) (bool, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		handleAuthErr := func() {
			utils.UnsetCookie(w, utils.Cookies.Login)
			rt := utils.RequestType(ctx)

			switch rt {
			case constants.RequestData, constants.RequestJSON:
				writeError(ctx, w, perr("failed to parse cookie, please clear your cookies and try again", http.StatusUnauthorized))
			default:
				utils.LogCtx(ctx).Panic("request type not set")
			}

		}

		var secret string
		var err error
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// try bearer token
			// split the header at the space character
			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
				handleAuthErr()
				return
			}
			decodedBytes, err := base64.StdEncoding.DecodeString(authHeaderParts[1])
			if err != nil {
				handleAuthErr()
				return
			}
			var tokenMap map[string]string
			err = json.Unmarshal(decodedBytes, &tokenMap)
			if err != nil {
				handleAuthErr()
				return
			}
			token, err := service.ParseAuthToken(tokenMap)
			if err != nil {
				handleAuthErr()
				return
			}
			secret = token.Secret
		} else {
			// try cookie
			secret, err = a.GetSecretFromCookie(ctx, r)
			if err != nil {
				handleAuthErr()
				return
			}
		}

		authInfo, ok, err := a.Service.GetSessionAuthInfo(ctx, secret)
		if err != nil {
			handleAuthErr()
			return
		}
		if !ok {
			handleAuthErr()
			return
		}

		if len(authorizers) == 0 {
			r = r.WithContext(context.WithValue(ctx, utils.CtxKeys.UserID, authInfo.UID))
			next(w, r)
			return
		}

		allOk := true

		for _, authorizer := range authorizers {
			ok, err := authorizer(r, authInfo.UID)
			if err != nil {
				utils.LogCtx(ctx).Error(err)
				writeError(ctx, w, perr("failed to verify authority", http.StatusInternalServerError))
				return
			}
			if !ok {
				allOk = false
				break
			}
		}

		if allOk {
			r = r.WithContext(context.WithValue(ctx, utils.CtxKeys.UserID, authInfo.UID))
			next(w, r)
			return
		}

		utils.LogCtx(ctx).Debug("unauthorized attempt")
		writeError(ctx, w, perr("you do not have the proper authorization to access this page", http.StatusUnauthorized))
	}
}

func (a *App) RequestJSON(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r.WithContext(context.WithValue(r.Context(), utils.CtxKeys.RequestType, constants.RequestJSON)))
	}
}

func (a *App) RequestData(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r.WithContext(context.WithValue(r.Context(), utils.CtxKeys.RequestType, constants.RequestData)))
	}
}
