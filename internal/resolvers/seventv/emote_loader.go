package seventv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Chatterino/api/internal/logger"
	"github.com/Chatterino/api/pkg/cache"
	"github.com/Chatterino/api/pkg/config"
	"github.com/Chatterino/api/pkg/resolver"
	"github.com/Chatterino/api/pkg/utils"
)

type EmoteLoader struct {
	apiURL  string
	baseURL string
}

func (l *EmoteLoader) Load(ctx context.Context, emoteHash string, r *http.Request) (*resolver.Response, time.Duration, error) {
	log := logger.FromContext(ctx)

	log.Debugw("[SevenTV] Get emote",
		"emoteHash", emoteHash,
	)

	queryMap := map[string]interface{}{
		"query": `
query fetchEmote($id: String!) {
	emote(id: $id) {
		visibility
		id
		name
		owner {
			id
			display_name
		}
	}
}`,
		"variables": map[string]string{
			"id": emoteHash,
		},
	}

	queryBytes, err := json.Marshal(queryMap)
	if err != nil {
		return resolver.Errorf("SevenTV API request marshal error: %s", err)
	}

	// Execute SevenTV API request
	resp, err := resolver.RequestPOST(l.apiURL, string(queryBytes))
	if err != nil {
		return resolver.Errorf("SevenTV API request error: %s", err)
	}
	defer resp.Body.Close()

	// Error out if the emote wasn't found or something else went wrong with the request
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		return emoteNotFoundResponse, cache.NoSpecialDur, nil
	}

	var jsonResponse EmoteAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		return resolver.Errorf("SevenTV API response decode error: %s", err)
	}

	// API returns Data.Emote as null if the emote wasn't found
	if jsonResponse.Data.Emote == nil {
		return emoteNotFoundResponse, cache.NoSpecialDur, nil
	}

	// Determine type of the emote based on visibility flags
	visibility := jsonResponse.Data.Emote.Visibility
	var emoteType []string

	if utils.HasBits(visibility, EmoteVisibilityGlobal) {
		emoteType = append(emoteType, "Global")
	}

	if utils.HasBits(visibility, EmoteVisibilityPrivate) {
		emoteType = append(emoteType, "Private")
	}

	// Default to Shared emote
	if len(emoteType) == 0 {
		emoteType = append(emoteType, "Shared")
	}

	// Build tooltip data from the API response
	data := TooltipData{
		Code:     jsonResponse.Data.Emote.Name,
		Type:     strings.Join(emoteType, " "),
		Uploader: jsonResponse.Data.Emote.Owner.DisplayName,
		Unlisted: utils.HasBits(visibility, EmoteVisibilityHidden),
	}

	// Build a tooltip using the tooltip template (see tooltipTemplate) with the data we massaged above
	var tooltip bytes.Buffer
	if err := seventvEmoteTemplate.Execute(&tooltip, data); err != nil {
		return resolver.Errorf("SevenTV emote template error: %s", err)
	}

	// Success
	successTooltip := &resolver.Response{
		Status:    http.StatusOK,
		Tooltip:   url.PathEscape(tooltip.String()),
		Thumbnail: utils.FormatThumbnailURL(l.baseURL, r, fmt.Sprintf(thumbnailFormat, emoteHash)),
		Link:      fmt.Sprintf("https://7tv.app/emotes/%s", emoteHash),
	}

	// Hide thumbnail for unlisted or hidden emotes pajaS
	if data.Unlisted {
		successTooltip.Thumbnail = ""
	}

	return successTooltip, cache.NoSpecialDur, nil
}

func NewEmoteLoader(cfg config.APIConfig, apiURL *url.URL) *EmoteLoader {
	return &EmoteLoader{
		apiURL:  apiURL.String(),
		baseURL: cfg.BaseURL,
	}
}
