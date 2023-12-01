package transport

import (
	"fmt"
	"mime"
	"net/http"
	"os"

	"github.com/FlashpointProject/CommunityWebsite/constants"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ResourceKeyUserID     = "user-id"
	ResourceKeyPlaylistID = "playlist-id"
	ResourceKeyGameID     = "game-id"
	ReosurceKeyUsername   = "username"
)

func (a *App) ServeRouter(l *logrus.Entry, srv *http.Server, router *mux.Router) {
	isStaff := func(r *http.Request, uid string) (bool, error) {
		return a.UserHasAnyRole(r, uid, constants.StaffRoles())
	}

	// Auth

	router.Handle("/auth/callback",
		http.HandlerFunc(a.RequestData(a.OAuthCallback))).
		Methods("GET")

	router.Handle("/auth",
		http.HandlerFunc(a.RequestData(a.OAuthLogin))).
		Methods("GET")

	// Profile

	router.Handle("/api/profile",
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(a.GetProfile)))).
		Methods("GET")

	router.Handle(fmt.Sprintf("/api/profile/{%s}", constants.ResourceKeyUserID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(a.GetUserProfile)))).
		Methods("GET")

	// Playlist

	router.Handle("/api/playlists",
		http.HandlerFunc(a.RequestJSON(a.SearchPlaylists))).
		Methods("GET")

	router.Handle("/api/playlists",
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(a.SubmitPlaylist)))).
		Methods("POST")

	router.Handle(fmt.Sprintf("/api/playlist/{%s}", constants.ResourceKeyPlaylistID),
		http.HandlerFunc(a.RequestJSON(a.GetPlaylist))).
		Methods("GET")

	router.Handle(fmt.Sprintf("/api/playlist/{%s}/preview", constants.ResourceKeyPlaylistID),
		http.HandlerFunc(a.RequestJSON(a.GetPlaylistPreview))).
		Methods("GET")

	router.Handle(fmt.Sprintf("/api/playlist/{%s}/download", constants.ResourceKeyPlaylistID),
		http.HandlerFunc(a.RequestJSON(a.DownloadPlaylist))).
		Methods("GET")

	router.Handle(fmt.Sprintf("/api/playlist/{%s}", constants.ResourceKeyPlaylistID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(a.UpdatePlaylist)))).
		Methods("PUT", "POST")

	router.Handle(fmt.Sprintf("/api/playlist/{%s}", constants.ResourceKeyPlaylistID),
		http.HandlerFunc(a.RequestJSON(a.UserAuthMux(a.DeletePlaylist)))).
		Methods("DELETE")

	// Games

	router.Handle(fmt.Sprintf("/api/game/{%s}", constants.ResourceKeyGameID),
		http.HandlerFunc(a.RequestJSON(a.GetGame))).
		Methods("GET")

	// GOTD

	router.Handle("/api/gotd/suggestions",
		http.HandlerFunc(a.RequestJSON(a.SearchGotdSuggestions))).
		Methods("GET")

	// News

	router.Handle(fmt.Sprintf("/api/post/{%s}", constants.ResourceKeyPostID),
		http.HandlerFunc(a.RequestJSON(a.GetNewsPost))).
		Methods("GET")

	router.Handle("/api/posts",
		http.HandlerFunc(a.RequestJSON(a.SearchNewsPosts))).
		Methods("GET")

	f := a.UserAuthMux(a.SubmitNewsPost, isStaff)

	router.Handle("/api/posts",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	// Content Reports

	f = a.UserAuthMux(a.SearchContentReports, isStaff)

	router.Handle("/api/reports",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("GET")

	f = a.UserAuthMux(a.SubmitContentReport)

	router.Handle("/api/reports",
		http.HandlerFunc(a.RequestJSON(f))).
		Methods("POST")

	// Roles

	router.Handle("/api/roles",
		http.HandlerFunc(a.RequestJSON(a.GetRoles))).
		Methods("GET")

	// Filter Groups

	router.Handle("/api/filter-groups",
		http.HandlerFunc(a.RequestJSON(a.GetFilterGroups))).
		Methods("GET")

	// Single Page Application (SPA)

	mime.AddExtensionType(".js", "application/javascript") // Windows being weird?
	router.PathPrefix("/").Handler(http.StripPrefix("/", &SinglePageAppHandler{
		Path: "./web/dist/index.html",
		Root: "/",
	}))

	err := srv.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}
}

type SinglePageAppHandler struct {
	Path string
	Root string
}

func (spa *SinglePageAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if the requested file exists
	path := "./web/dist/" + r.URL.Path
	if r.URL.Path == "index.html" {
		http.Redirect(w, r, spa.Root, http.StatusFound)
		return
	}
	stats, err := os.Stat(path)
	if os.IsNotExist(err) {
		// If not, serve index.html
		http.ServeFile(w, r, spa.Path)
	} else {
		if stats.IsDir() {
			// If it's a directory, serve index.html
			http.ServeFile(w, r, spa.Path)
			return
		}
		// If it exists, serve the file
		http.ServeFile(w, r, path)
	}
}
