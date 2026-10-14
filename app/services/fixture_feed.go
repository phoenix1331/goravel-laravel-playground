// Package services holds plain Go services with no framework lookups.
package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/goravel/framework/contracts/http/client"
)

// FeedResult is one result as published by the league.
type FeedResult struct {
	Home      string `json:"home"`
	Away      string `json:"away"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

// FixtureFeed is anywhere results can be read from: a JSON file or the
// league's HTTP API. Both types below satisfy it without saying so; Go has
// no "implements" keyword.
type FixtureFeed interface {
	Results() ([]FeedResult, error)
}

// FileFeed reads results from a JSON file.
type FileFeed struct {
	Path string
}

func (f FileFeed) Results() ([]FeedResult, error) {
	contents, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, fmt.Errorf("read fixture feed: %w", err)
	}

	var results []FeedResult
	if err := json.Unmarshal(contents, &results); err != nil {
		return nil, fmt.Errorf("parse fixture feed %s: %w", f.Path, err)
	}

	return results, nil
}

// HttpFeed reads results from the league's API. It is given a ready-made
// request builder, so it never reads config itself and tests can fake it.
//
// Laravel:
//
//	Http::baseUrl($url)->acceptJson()->get('/results')->throw()->json();
//
// Goravel: very close, fakes included. A 4xx or 5xx is not an error, so
// check response.Failed() where Laravel would call ->throw().
type HttpFeed struct {
	Request client.Request
}

func (f HttpFeed) Results() ([]FeedResult, error) {
	response, err := f.Request.AcceptJSON().Get("/results")
	if err != nil {
		// The request never completed: DNS, connection refused, timeout
		return nil, fmt.Errorf("fetch fixture feed: %w", err)
	}
	if response.Failed() {
		// The request completed with a 4xx or 5xx. Unlike Guzzle, the client
		// does not treat that as an error, so check the status yourself.
		return nil, fmt.Errorf("fetch fixture feed: status %d", response.Status())
	}

	var results []FeedResult
	if err := response.Bind(&results); err != nil {
		return nil, fmt.Errorf("parse fixture feed: %w", err)
	}

	return results, nil
}
