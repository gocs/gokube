package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gocs/gokube/internal/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	ID                     string `json:"id"`
	Title                  string `json:"title"`
	VideoOwnerChannelTitle string `json:"video_owner_channel_title"`
	Thumbnail              string `json:"thumbnail"`
	URL                    string `json:"url"`
	ViewCount              int64  `json:"view_count"`
}

type Playlist struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Thumbnail   string  `json:"thumbnail"`
	Videos      []Video `json:"videos"`
}

type YoutubeCtrl struct {
	service *youtube.Service
}

func NewYoutubeCtrl(ctx context.Context) (*YoutubeCtrl, error) {
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(config.YoutubeAPIKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube service: %v", err)
	}
	return &YoutubeCtrl{
		service: youtubeService,
	}, nil
}

func (yc *YoutubeCtrl) Playlist(w http.ResponseWriter, r *http.Request) {
	playlistID := r.URL.Query().Get("playlist_id")
	if playlistID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// First, get playlist details
	playlistCall := yc.service.Playlists.List([]string{"snippet"}).Id(playlistID)
	playlistResponse, err := playlistCall.Do()
	if err != nil {
		log.Printf("Error getting playlist details: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(playlistResponse.Items) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	playlistInfo := playlistResponse.Items[0]
	playlist := Playlist{
		ID:          playlistID,
		Name:        playlistInfo.Snippet.Title,
		Description: playlistInfo.Snippet.Description,
		Thumbnail:   playlistInfo.Snippet.Thumbnails.Default.Url,
		Videos:      []Video{},
	}

	// Get all videos from the playlist (handle pagination)
	var allVideos []Video
	nextPageToken := ""

	for {
		call := yc.service.PlaylistItems.List([]string{"snippet"}).PlaylistId(playlistID).MaxResults(50)
		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		response, err := call.Do()
		if err != nil {
			log.Printf("Error getting playlist items: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Add videos from this page
		wg := sync.WaitGroup{}
		for _, item := range response.Items {
			wg.Add(1)
			go func(item *youtube.PlaylistItem) {
				defer wg.Done()
				video, err := yc.Video(item.Snippet.ResourceId.VideoId)
				if err != nil {
					log.Printf("Error getting video details: %v", err)
					return
				}
				allVideos = append(allVideos, video)
			}(item)
		}
		wg.Wait()

		// Check if there are more pages
		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	playlist.Videos = allVideos
	log.Printf("Retrieved %d videos from playlist", len(allVideos))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlist)
}

func (yc *YoutubeCtrl) Video(videoID string) (Video, error) {

	// Get video details including statistics
	call := yc.service.Videos.List([]string{"snippet", "statistics"}).Id(videoID)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error getting video details: %v", err)
		return Video{}, fmt.Errorf("error getting video details: %v", err)
	}

	if len(response.Items) == 0 {
		return Video{}, fmt.Errorf("video not found")
	}

	videoInfo := response.Items[0]
	video := Video{
		ID:                     videoID,
		Title:                  videoInfo.Snippet.Title,
		VideoOwnerChannelTitle: videoInfo.Snippet.ChannelTitle,
		Thumbnail:              videoInfo.Snippet.Thumbnails.Default.Url,
		URL:                    fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		ViewCount:              int64(videoInfo.Statistics.ViewCount),
	}

	return video, nil
}
