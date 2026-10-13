<?php

namespace App\Providers;

use App\Services\LeagueTable;
use App\Services\PointsLeagueTable;
use Illuminate\Support\ServiceProvider;

class LeagueServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        // Bound by interface name: anything type-hinting LeagueTable gets this
        $this->app->singleton(LeagueTable::class, fn ($app) => new PointsLeagueTable(
            config('club.points.win'),
            config('club.points.draw'),
            config('club.points.loss'),
        ));
    }

    public function boot(): void
    {
        //
    }
}
