<?php

namespace App\Listeners;

use App\Events\FixtureResultRecorded;
use Illuminate\Support\Facades\Cache;

class ForgetLeagueTable
{
    public function handle(FixtureResultRecorded $event): void
    {
        Cache::forget('league_table');
    }
}
