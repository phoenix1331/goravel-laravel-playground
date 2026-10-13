<?php

public function test_count_creates_many_rows(): void
{
    $clubs = Club::factory()->count(3)->create();

    $this->assertCount(3, $clubs);
}

public function test_attributes_override_the_definition(): void
{
    $player = Player::factory()->for(Club::factory())->create(['squad_number' => 7]);

    $this->assertSame(7, $player->squad_number);
}

public function test_make_does_not_save(): void
{
    $club = Club::factory()->make();

    $this->assertFalse($club->exists);
}

public function test_database_seeder_builds_a_league(): void
{
    $this->seed(DatabaseSeeder::class);

    $this->assertDatabaseCount('clubs', 6);
    $this->assertDatabaseCount('players', 96);
    $this->assertDatabaseCount('fixtures', 15);
    $this->assertDatabaseCount('users', 6);
}
