<?php

namespace Tests\Unit;

use App\Services\LeagueTable;
use App\Services\PointsLeagueTable;
use App\Services\Result;
use PHPUnit\Framework\Attributes\DataProvider;
use Tests\TestCase;

class LeagueTableTest extends TestCase
{
    public static function rankings(): array
    {
        return [
            'a win beats a loss' => [[new Result(1, 2, 0, 2)], [2, 1]],
            'level on points, goal difference decides' => [
                [new Result(1, 3, 1, 0), new Result(2, 3, 4, 0)],
                [2, 1, 3],
            ],
            'a draw is level on everything, so lower club ID first' => [[new Result(2, 1, 1, 1)], [1, 2]],
        ];
    }

    #[DataProvider('rankings')]
    public function test_points_league_table_ranking(array $results, array $want): void
    {
        $table = new PointsLeagueTable(3, 1, 0);

        $this->assertSame($want, array_column($table->standings($results), 'club_id'));
    }

    public function test_container_resolves_league_table_singleton(): void
    {
        $this->assertSame(app(LeagueTable::class), app(LeagueTable::class));
    }
}
