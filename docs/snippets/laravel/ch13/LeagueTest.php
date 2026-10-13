<?php

public function test_table_ranks_clubs_from_finished_fixtures(): void
{
    $this->getJson('/api/league/table')
        ->assertOk()
        ->assertJsonCount(4, 'data')
        ->assertJsonPath('data.0.club', 'Ashford')
        ->assertJsonPath('data.0.points', 3)
        ->assertJsonPath('data.0.played', 1)
        ->assertJsonPath('data.1.club', 'Crewe')
        ->assertJsonPath('data.2.club', 'Barnsley')
        ->assertJsonPath('data.3.club', 'Dorking');
}

public function test_goals_totals_home_and_away_goals(): void
{
    $this->getJson('/api/league/goals')
        ->assertOk()
        ->assertJsonCount(3, 'data')
        ->assertJsonPath('data.0.club', 'Ashford')
        ->assertJsonPath('data.0.goals', 3);
}
