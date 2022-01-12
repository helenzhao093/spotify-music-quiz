package main


import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var (
	spotifyOauthConfig *oauth2.Config
	client             *http.Client
	oauthStateString   = "abc123"
	state              string
	code               string
	templates          = template.Must(template.ParseFiles("client/build/index.html"))
	trackInfo          track
	tracksInfo         tracks
	playlistsInfo      playlists
	token              string
)

const (
	ScopeUserLibraryRead     = "user-library-read"
	ScopePlaylistReadPrivate = "playlist-read-private"
	ScopeUserReadEmail       = "user-read-email"
	ScopeStreaming           = "streaming"
)

type Page struct {
	Title string
	Body  []byte
}

type externalUrls struct {
	spotify string `json:"spotify"`
}

type artist struct {
	ExternalUrls externalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	URI          string       `json:"uri"`
	Type         string       `json:"type"`
}

type image struct {
	height int    `json:"height"`
	url    string `json:"url"`
	width  int    `json:"width"`
}

type album struct {
	AlbumType           string       `json:"album_type"`
	Artists             []artist     `json:"artists"`
	AvailableMarkets    []string     `json:"available_markets"`
	ExternalUrls        externalUrls `json:"external_urls"`
	Href                string       `json:"href"`
	ID                  string       `json:"id"`
	Images              []image      `json:"images"`
	Name                string       `json:"name"`
	ReleaseDate         string       `json:"release_date"`
	ReleaseDatPrecision string       `json:"release_date_precision"`
	URI                 string       `json:"uri"`
	TotalTracks         int          `json:"total_tracks"`
	Type                string       `json:"type"`
}

type externaIds struct {
	Isrc string `json:"isrc"`
}

type track struct {
	Album            album        `json:"album"`
	Artists          []artist     `json:"artists"`
	AvailableMarkets []string     `json:"available_markets"`
	URI              string       `json:"uri"`
	Type             string       `json:"type"`
	PreviewURL       string       `json:"preview_url"`
	Name             string       `json:"name"`
	DiscNumber       int          `json:"disc_number"`
	DurationMS       int          `json:"duration_ms"`
	Episode          bool         `json:"episode"`
	Explicit         bool         `json:"explicit"`
	ExternalIds      externaIds   `json:"external_ids"`
	ExternalUrls     externalUrls `json:"external_urls"`
	Href             string       `json:"href"`
	ID               string       `json:"id"`
	IsLocal          bool         `json:"is_local"`
	IsPlayable       bool         `json:"is_playable"`
	Popularity       int          `json:"popularity"`
	TrackNumber      int          `json:"track_number"`
	Track            bool         `json:"track"`
}

type tracks struct {
	Href     string `json:"href"`
	Items    []item `json:"items"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
}

type videoThumbnail struct {
	URL string `json:"video_thumbnail"`
}

type item struct {
	AddedAt        string         `json:"added_at"`
	AddedBy        user           `json:"added_by"`
	IsLocal        bool           `json:"is_local"`
	PrimaryColor   string         `json:"primary_color"`
	Track          track          `json:"track"`
	VideoThumbnail videoThumbnail `json:"video_thumbnail"`
}

type previewUrl struct {
	URL string
}

type trackClient struct {
	ID         string
	Name       string
	PreviewURL string
}

type followers struct {
	Href  string `json:"href"`
	Total int32  `json:"total"`
}

type user struct {
	ExternalUrls externalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type playlist struct {
	Collaborative bool         `json:"collaborative"`
	Description   string       `json:"description"`
	ExternalUrls  externalUrls `json:"external_urls"`
	Followers     followers    `json:"followers"`
	Href          string       `json:"href"`
	ID            string       `json:"id"`
	Images        []image      `json:"images"`
	Name          string       `json:"name"`
	Owner         user         `json:"owner"`
	PrimaryColor  string       `json:"primary_color"`
	Public        string       `json:"public"`
	SnapshotID    string       `json:"snapshot_id"`
	Tracks        tracks       `json:"tracks"`
	Type          string       `json:"type"`
	URI           string       `json:"uri"`
}

type playlists struct {
	Href      string     `json:"href"`
	Playlists []playlist `json:"items"`
	Limit     int        `json:"int"`
	Next      string     `json:"next"`
	Offset    int        `json:"offset"`
	Previous  string     `json:"previous"`
	Total     int        `json:"total"`
}

type playlistClient struct {
	ID   string
	Name string
}

type playlistsClient struct {
	Playlists []playlistClient
}

func loadPage(title string) *Page {
	filename := title + ".txt"
	body, _ := ioutil.ReadFile(filename)
	return &Page{Title: title, Body: body}
}

func init() {
	os.Setenv("SPOTIFY_ID", "05626d01cfe2402abf20a9b399dfb69b")
	os.Setenv("SPOTIFY_SECRET", "173346750aa44f17acfd5190d3b35550")
	spotifyOauthConfig = &oauth2.Config{
		ClientID:     "05626d01cfe2402abf20a9b399dfb69b", //os.Getenv("SPOTIFY_ID"),
		ClientSecret: "173346750aa44f17acfd5190d3b35550", //os.Getenv("SPOTIFY_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{ScopeUserLibraryRead, ScopePlaylistReadPrivate, ScopeUserReadEmail, ScopeStreaming},
		Endpoint:     spotify.Endpoint,
	}
	fmt.Println(os.Getenv("SPOTIFY_ID"))
}

func handleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := spotifyOauthConfig.AuthCodeURL(oauthStateString)
	//c.Redirect(http.StatusTemporaryRedirect, url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

/*func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.Execute(w, tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}*/

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<body>
	<a href="/login">Spotify Log In</a>
</body>
</html>`
	fmt.Fprintf(w, htmlIndex)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	err := templates.Execute(w, "client/build/index.html")
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
}

func handleSpotifyCallback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	state, code = r.FormValue("state"), r.FormValue("code")

	fmt.Println(state)
	fmt.Println(code)
	configClient(state, code)
	http.Redirect(w, r, "/index", http.StatusTemporaryRedirect)
}

func configClient(state string, code string) {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	token, _ := spotifyOauthConfig.Exchange(ctx, code)
	fmt.Println(token)
	client = spotifyOauthConfig.Client(ctx, token)
}

func handleGetTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackId := vars["trackId"]
	fmt.Println(vars["trackId"])
	url := "https://api.spotify.com/v1/tracks/" + trackId + "?market=US"
	response, _ := client.Get(url)
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(contents, &trackInfo)
	fmt.Printf("preview url: %s,", trackInfo.PreviewURL)

	preview := previewUrl{trackInfo.PreviewURL}
	b, _ := json.Marshal(preview)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func handleGetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	url := "https://api.spotify.com/v1/me/playlists?limit=10"
	response, _ := client.Get(url)
	defer response.Body.Close()
	playlistsByte, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(playlistsByte))
	_ = json.Unmarshal(playlistsByte, &playlistsInfo)
	fmt.Printf("%+v\n", playlistsInfo)
	var playlistsClient []playlistClient
	for _, playlist := range playlistsInfo.Playlists {
		playlistClient := playlistClient{playlist.ID, playlist.Name}
		playlistsClient = append(playlistsClient, playlistClient)
	}
	b, _ := json.Marshal(playlistsClient)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func handleGetTracksFromPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistId := vars["playlistId"]
	fmt.Println(vars["playlistId"])
	url := "https://api.spotify.com/v1/playlists/" + playlistId + "/tracks?market=US&limit=50"
	fmt.Println(url)
	response, _ := client.Get(url)
	defer response.Body.Close()
	tracksResponse, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(tracksResponse))
	_ = json.Unmarshal(tracksResponse, &tracksInfo)
	var tracksClient []trackClient
	for _, item := range tracksInfo.Items {
		trackClient := trackClient{item.Track.ID, item.Track.Name, item.Track.PreviewURL}
		tracksClient = append(tracksClient, trackClient)
	}
	b, _ := json.Marshal(tracksClient)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	r := mux.NewRouter()
	buildPath := path.Clean("client/build")
	buildURL := fmt.Sprintf("/%s/", buildPath)
	r.Handle(buildURL, http.StripPrefix(buildURL, http.FileServer(http.Dir(buildPath))))

	staticHandler := http.FileServer(http.Dir("client/build/static"))
	r.PathPrefix("/static/").Handler(staticHandler)

	//http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("views"))))
	r.HandleFunc("/home", handleMain)
	r.HandleFunc("/login", handleSpotifyLogin)
	r.HandleFunc("/callback", handleSpotifyCallback)
	r.HandleFunc("/index", handleIndex)
	r.HandleFunc("/getTrack/{trackId}", handleGetTrack)
	r.HandleFunc("/getPlaylists", handleGetUserPlaylists)
	r.HandleFunc("/getTracksFromPlaylist/{playlistId}", handleGetTracksFromPlaylist)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
