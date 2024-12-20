package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var body io.Reader
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, body)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")

    var client http.Client
    res, err := client.Do(req)
    if err != nil {
		return &RSSFeed{}, err
    }
    defer res.Body.Close()

    data, err := io.ReadAll(res.Body)
    if err != nil {
		return &RSSFeed{}, err
    }

    var rssFeed RSSFeed
    err = xml.Unmarshal(data, &rssFeed)
    if err != nil {
		return &RSSFeed{}, err
    }

    rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
    rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
    for i := range rssFeed.Channel.Item {
        rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
        rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
    }

	return &rssFeed, nil
}
