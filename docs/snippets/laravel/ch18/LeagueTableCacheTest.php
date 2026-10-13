<?php

public function test_table_is_cached_after_first_request(): void
{
    $this->getJson('/api/league/table');

    $this->assertTrue(Cache::has('league_table'));
}

public function test_recording_a_result_clears_the_cache(): void
{
    $this->getJson('/api/league/table');

    app(RecordResult::class)($this->fixture, 1, 0);

    $this->assertFalse(Cache::has('league_table'));
    $this->getJson('/api/league/table')
        ->assertJsonPath('data.0.club', 'Ashford')
        ->assertJsonPath('data.0.points', 3);
}
