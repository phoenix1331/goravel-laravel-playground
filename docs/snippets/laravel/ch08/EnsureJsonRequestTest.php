<?php

namespace Tests\Feature;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class EnsureJsonRequestTest extends TestCase
{
    use RefreshDatabase;

    public function test_rejects_non_json_writes(): void
    {
        $this->post('/api/clubs', ['name' => 'Form FC', 'ground' => 'Park'])
            ->assertStatus(415);

        $this->assertDatabaseCount('clubs', 0);
    }

    public function test_allows_json_writes(): void
    {
        $this->postJson('/api/clubs', ['name' => 'Json FC', 'ground' => 'Park'])
            ->assertCreated();
    }

    public function test_reads_are_unaffected(): void
    {
        $this->getJson('/api/clubs')->assertOk();
    }
}
