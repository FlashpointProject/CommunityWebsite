package transport

import (
	"net/http"
)

type commCtxKey string

type commCtxKeys struct {
	FPFSS commCtxKey
}

var CommCtxKeys = commCtxKeys{
	FPFSS: "fpfss",
}

func (a *App) UserHasAnyRole(r *http.Request, uid string, roles []string) (bool, error) {
	ctx := r.Context()

	// @TODO Cache

	user, err := a.Service.GetUser(ctx, uid)
	if err != nil {
		return false, err
	}

	isAuthorized := HasAnyRole(user.Roles, roles)
	if !isAuthorized {
		return false, nil
	}

	return true, nil
}

func HasAnyRole(has, needs []string) bool {
	for _, role := range has {
		for _, neededRole := range needs {
			if role == neededRole {
				return true
			}
		}
	}
	return false
}
