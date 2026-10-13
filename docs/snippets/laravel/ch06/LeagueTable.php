<?php

namespace App\Facades;

use Illuminate\Support\Facades\Facade;

/**
 * @method static array standings(array $results)
 */
class LeagueTable extends Facade
{
    // Every static call is forwarded, via __callStatic, to whatever the
    // container holds under this key
    protected static function getFacadeAccessor(): string
    {
        return \App\Services\LeagueTable::class;
    }
}
