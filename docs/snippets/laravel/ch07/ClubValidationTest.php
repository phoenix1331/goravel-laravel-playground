<?php

namespace Tests\Feature;

use App\Models\Club;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class ClubValidationTest extends TestCase
{
    use RefreshDatabase;

    public function test_required_fields(): void
    {
        $this->postJson('/api/clubs', [])
            ->assertUnprocessable()
            ->assertJsonValidationErrors(['name', 'ground'])
            ->assertJsonMissingValidationErrors(['founded_year']);
    }

    public function test_whitespace_is_trimmed_before_required_check(): void
    {
        $this->postJson('/api/clubs', ['name' => '   ', 'ground' => 'Park'])
            ->assertJsonValidationErrors(['name']);
    }

    public function test_founded_year_range(): void
    {
        $this->postJson('/api/clubs', ['name' => 'Old FC', 'ground' => 'Park', 'founded_year' => 1800])
            ->assertJsonValidationErrors([
                'founded_year' => 'The oldest football club was founded in 1857.',
            ]);
    }

    public function test_name_must_be_unique(): void
    {
        Club::create(['name' => 'Taken FC', 'ground' => 'Park']);

        $this->postJson('/api/clubs', ['name' => 'Taken FC', 'ground' => 'Park'])
            ->assertJsonValidationErrors(['name' => 'A club with this name already exists.']);
    }
}
