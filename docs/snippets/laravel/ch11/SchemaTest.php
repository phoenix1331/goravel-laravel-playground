<?php

public function test_schema_has_expected_columns_and_indexes(): void
{
    $this->assertTrue(Schema::hasColumns('clubs', ['id', 'name', 'ground', 'founded_year', 'short_name']));
    $this->assertTrue(Schema::hasIndex('fixtures', 'fixtures_kick_off_at_index'));
}

public function test_database_enforces_unique_squad_number_per_club(): void
{
    $club = Club::create(['name' => 'Ashford Athletic', 'ground' => 'The Rec']);
    $club->players()->create(['name' => 'A', 'position' => 'forward', 'squad_number' => 9]);

    $this->expectException(UniqueConstraintViolationException::class);
    $club->players()->create(['name' => 'B', 'position' => 'forward', 'squad_number' => 9]);
}

public function test_kick_off_is_stored_in_utc(): void
{
    $fixture = Fixture::create([/* ... */ 'kick_off_at' => '2026-10-10 15:00:00+01:00']);

    $this->assertSame('14:00', $fixture->fresh()->kick_off_at->utc()->format('H:i'));
}
