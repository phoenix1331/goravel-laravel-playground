// Laravel equivalent: app/Services/FixtureFeed.php - see docs/chapters/14-console-commands.md
package services

import (
	"encoding/json"
	"fmt"
	"os"
)

// FeedResult is one result as published by the league.
type FeedResult struct {
	Home      string `json:"home"`
	Away      string `json:"away"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

// #region feed
// FixtureFeed is anywhere results can be read from. Chapter 21 adds an HTTP
// implementation next to this file-based one.
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

// #endregion feed
