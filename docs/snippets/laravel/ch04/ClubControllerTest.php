<?php

namespace Tests\Feature;

use App\Models\Club;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class ClubControllerTest extends TestCase
{
    use RefreshDatabase;

    public function test_index_lists_clubs_by_name(): void
    {
        Club::create(['name' => 'Riverside Rovers', 'ground' => 'Riverside Park']);
        Club::create(['name' => 'Ashford Athletic', 'ground' => 'Ashford Park']);

        $this->getJson('/api/clubs')
            ->assertOk()
            ->assertJsonCount(2, 'data')
            ->assertJsonPath('data.0.name', 'Ashford Athletic')
            ->assertJsonPath('data.1.name', 'Riverside Rovers');
    }

    public function test_show_returns_404_for_missing_club(): void
    {
        $this->getJson('/api/clubs/999')->assertNotFound();
    }

    public function test_store_creates_club(): void
    {
        $this->postJson('/api/clubs', [
            'name' => 'Hillside United',
            'ground' => 'Hillside Rec',
            'founded_year' => 1987,
        ])->assertCreated();

        $this->assertDatabaseHas('clubs', [
            'name' => 'Hillside United',
            'ground' => 'Hillside Rec',
            'founded_year' => 1987,
        ]);
    }
}
