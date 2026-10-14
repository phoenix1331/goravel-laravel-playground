// This package uses Goravel's generated mocks. mock.Factory() swaps the whole
// application for a mock, process-wide, so it must not share a package with
// tests that boot the real app. go test runs each package in its own process,
// and this one never imports goravel/tests, so nothing real is booted here.
package mocked

import (
	"testing"

	"github.com/goravel/framework/testing/mock"

	"goravel/app/listeners"
	"goravel/app/services"
)

func TestForgetLeagueTableForgetsTheCachedTable(t *testing.T) {
	// Like Cache::shouldReceive('forget')->once()->with('league_table')
	mockCache := mock.Factory().Cache()
	mockCache.EXPECT().Forget(services.LeagueTableCacheKey).Return(true).Once()

	err := (&listeners.ForgetLeagueTable{}).Handle(uint(1))

	if err != nil {
		t.Fatalf("Handle returned %v", err)
	}
	mockCache.AssertExpectations(t)
}
