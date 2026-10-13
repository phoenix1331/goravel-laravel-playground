<?php

namespace Tests\Unit;

use App\Facades\LeagueTable;
use App\Services\Result;
use Tests\TestCase;

class LeagueTableFacadeTest extends TestCase
{
    public function test_league_table_facade_resolves_binding(): void
    {
        $standings = LeagueTable::standings([new Result(1, 2, 2, 0)]);

        $this->assertCount(2, $standings);
        $this->assertSame(1, $standings[0]->clubId);
        $this->assertSame(3, $standings[0]->points);
    }

    public function test_league_table_facade_can_be_swapped(): void
    {
        // Facades can be mocked directly; Laravel restores them after each test
        LeagueTable::shouldReceive('standings')
            ->once()
            ->andReturn([['club_id' => 99, 'points' => 100]]);

        $this->assertSame(99, LeagueTable::standings([])[0]['club_id']);
    }
}
