package mux

import (
	"fmt"
	"log"
	"os"

	muxgo "github.com/muxinc/mux-go/v5"
)

var (
	client *muxgo.APIClient
)

func init() {
	muxTokenID := os.Getenv("MUX_TOKEN_ID")
	muxTokenSecret := os.Getenv("MUX_TOKEN_SECRET")

	if muxTokenID == "" || muxTokenSecret == "" {
		log.Fatalf("Mux token ID or secret are not set in environment")
	}

	client = muxgo.NewAPIClient(
		muxgo.NewConfiguration(
			muxgo.WithBasicAuth(
				muxTokenID,
				muxTokenSecret,
			),
		),
	)

	log.Println("Mux client initialized")
}

func CreateAsset(s3Url string) (string, error) {
	asset, err := client.AssetsApi.CreateAsset(muxgo.CreateAssetRequest{
		Input: []muxgo.InputSettings{
			{
				Url: s3Url,
				// TextType:       "subtitles",
				// Type:           "text",
				// ClosedCaptions: false,
			},
		},
		PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC},
	})
	if err != nil {
		return "", err
	}

	muxUrl := fmt.Sprintf("https://stream.mux.com/%s.m3u8", asset.Data.PlaybackIds[0].Id)

	return muxUrl, nil
}
