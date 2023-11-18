package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/FlashpointProject/CommunityWebsite/constants"
	"github.com/FlashpointProject/CommunityWebsite/service"
	"github.com/FlashpointProject/CommunityWebsite/utils"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// BoolString is a little hack to make handling tri-state bool in go templates trivial
func BoolString(b *bool) string {
	if b == nil {
		return "nil"
	}
	if *b {
		return "true"
	}
	return "false"
}

func (a *App) GetUserIDFromCookie(r *http.Request) (int64, error) {
	cookieMap, err := a.CC.GetSecureCookie(r, utils.Cookies.Login)
	if errors.Is(err, http.ErrNoCookie) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	token, err := service.ParseAuthToken(cookieMap)
	if err != nil {
		return 0, err
	}

	uid, err := strconv.ParseInt(token.UserID, 10, 64)
	if err != nil {
		return 0, err
	}

	return uid, nil
}

func (a *App) GetSecretFromCookie(ctx context.Context, r *http.Request) (string, error) {
	cookieMap, err := a.CC.GetSecureCookie(r, utils.Cookies.Login)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", err
		}
		utils.LogCtx(ctx).Error(err)
		return "", err
	}

	token, err := service.ParseAuthToken(cookieMap)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return "", err
	}

	return token.Secret, nil
}

func writeResponse(ctx context.Context, w http.ResponseWriter, data interface{}, status int) {
	requestType := utils.RequestType(ctx)
	if requestType == "" {
		utils.LogCtx(ctx).Panic("request type not set")
		return
	}

	switch requestType {
	case constants.RequestJSON, constants.RequestData:
		w.WriteHeader(status)
		if data != nil {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(data)
			if err != nil {
				utils.LogCtx(ctx).Error(err)
				if errors.Is(err, syscall.ECONNRESET) {
					return
				}
				if errors.Is(err, syscall.EPIPE) {
					return
				}
				writeError(ctx, w, err)
			}
		}
	default:
		utils.LogCtx(ctx).Panic("unsupported request type")
	}
}

func writeError(ctx context.Context, w http.ResponseWriter, err error) {
	ufe := &constants.PublicError{}
	if errors.As(err, ufe) {
		writeResponse(ctx, w, presp(ufe.Msg, ufe.Status), ufe.Status)
	} else {
		msg := http.StatusText(http.StatusInternalServerError)
		writeResponse(ctx, w, presp(msg, http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func perr(msg string, status int) error {
	return constants.PublicError{Msg: msg, Status: status}
}

func dberr(err error) error {
	return constants.DatabaseError{Err: err}
}

func presp(msg string, status int) constants.PublicResponse {
	return constants.PublicResponse{Msg: &msg, Status: status}
}

func capString(maxLen int, s *string) string {
	if s == nil {
		return "<nil>"
	}
	str := *s
	if len(str) <= 3 {
		return *s
	}
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen-3] + "..."
}

func milliTime(date time.Time) int64 {
	return date.UnixMilli()
}

func localeNum(ref *int64) string {
	p := message.NewPrinter(language.English)
	s := p.Sprintf("%d\n", *ref)
	return s
}

func isReturnURLValid(s string) bool {
	return len(s) > 0 && strings.HasPrefix(s, "/") && !strings.HasPrefix(s, "//")
}

func isReturnURLLocal(s string, host string) bool {
	if len(s) == 0 {
		return true
	}
	if strings.HasPrefix(s, "/") {
		return true
	}
	return strings.HasPrefix(s, host)
}

func SetCookie(w http.ResponseWriter, name string, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   maxAge,
	}
	http.SetCookie(w, cookie)
}
