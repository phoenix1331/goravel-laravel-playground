<?php

namespace Tests\Feature;

use Tests\TestCase;

class RoutingTest extends TestCase
{
    public function test_closure_route_inside_api_group(): void
    {
        $this->getJson('/api/league')
            ->assertOk()
            ->assertExactJson(['name' => 'Test League']);
    }

    public function test_named_handler_route(): void
    {
        $this->getJson('/api/league/points')
            ->assertOk()
            ->assertExactJson(['win' => 3, 'draw' => 1, 'loss' => 0]);
    }

    public function test_route_outside_group_is_not_prefixed(): void
    {
        $this->getJson('/league')->assertNotFound();
    }
}
