<?php

public function test_sync_reads_results_from_the_api(): void
{
    Http::preventStrayRequests();
    Http::fake([
        'https://feed.test/api/results' => Http::response([
            ['home' => 'Ashford', 'away' => 'Barnsley', 'home_score' => 2, 'away_score' => 0],
        ]),
    ]);

    $this->artisan('sync:fixtures')->assertSuccessful();

    $this->assertSame(2, $this->fixture->fresh()->home_score);
    Http::assertSent(fn ($request) => $request->hasHeader('Accept', 'application/json'));
}

public function test_server_error_fails_the_command_and_changes_nothing(): void
{
    Http::fake(['feed.test/*' => Http::response(status: 503)]);

    $this->artisan('sync:fixtures')->assertFailed();
}
