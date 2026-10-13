<?php

// app/Http/Controllers/LeagueController.php
public function table(): JsonResponse
{
    // Any serialisable value can be cached: Laravel serialize()s it for Redis
    $rows = Cache::remember('league_table', now()->addHour(), fn () => $this->buildLeagueTable());

    return response()->json(['data' => $rows]);
}
