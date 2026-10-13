<?php

namespace App\Services;

interface LeagueTable
{
    /** @param  array<int, Result>  $results */
    public function standings(array $results): array;
}

class PointsLeagueTable implements LeagueTable
{
    public function __construct(
        private int $win,
        private int $draw,
        private int $loss,
    ) {}

    public function standings(array $results): array
    {
        // ...
    }
}
