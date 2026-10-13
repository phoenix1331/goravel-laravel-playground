<?php

public function test_club_players_are_loaded_in_squad_order(): void
{
    $club = Club::create(['name' => 'Ashford Athletic', 'ground' => 'Ashford Park']);
    $club->players()->createMany([
        ['name' => 'Nine', 'position' => 'forward', 'squad_number' => 9],
        ['name' => 'One', 'position' => 'goalkeeper', 'squad_number' => 1],
    ]);

    $this->getJson("/api/clubs/{$club->id}/players")
        ->assertOk()
        ->assertJsonPath('data.players.0.name', 'One')
        ->assertJsonPath('data.players.1.name', 'Nine');
}

public function test_fixtures_eager_load_both_clubs(): void
{
    $home = Club::create(['name' => 'Home Town', 'ground' => 'Home Town Park']);
    $away = Club::create(['name' => 'Away Days', 'ground' => 'Away Days Park']);
    Fixture::create([
        'home_club_id' => $home->id, 'away_club_id' => $away->id,
        'kick_off_at' => '2026-11-07 15:00:00', 'venue' => 'Home Town Park', 'status' => 'scheduled',
    ]);

    $this->getJson('/api/fixtures')
        ->assertOk()
        ->assertJsonPath('data.0.home_club.name', 'Home Town')
        ->assertJsonPath('data.0.away_club.name', 'Away Days');
}
