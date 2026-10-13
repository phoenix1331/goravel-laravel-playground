<?php

namespace Tests\Feature;

use App\Models\Club;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class ClubControllerTest extends TestCase
{
    // Migrates once, then wraps each test in a transaction it rolls back
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        // runs before each test
    }

    public function test_index_lists_clubs_by_name(): void
    {
        Club::factory()->create(['name' => 'Riverside Rovers']);

        $this->getJson('/api/clubs')->assertOk()->assertJsonCount(1, 'data');
    }
}
