<?php

namespace App\Http\Controllers;

use App\Facades\LeagueTable;
use App\Services\Result;
use Illuminate\Support\Facades\DB;

class LeagueController extends Controller
{
    public function table(): JsonResponse
    {
        // Query builder: rows come back as stdClass objects
        $results = DB::table('fixtures')
            ->select('home_club_id', 'away_club_id', 'home_score', 'away_score')
            ->where('status', 'finished')
            ->get()
            ->map(fn ($row) => new Result(...(array) $row))
            ->all();

        $standings = LeagueTable::standings($results);

        // ... add positions, names and clubs that have not played
    }

    public function goals(): JsonResponse
    {
        // Raw SQL: ? placeholders on every driver, PDO handles the rest
        $rows = DB::select(<<<'SQL'
            SELECT clubs.id AS club_id, clubs.name AS club, SUM(scored.goals) AS goals
            FROM (
                SELECT home_club_id AS club_id, home_score AS goals FROM fixtures WHERE status = ?
                UNION ALL
                SELECT away_club_id AS club_id, away_score AS goals FROM fixtures WHERE status = ?
            ) AS scored
            JOIN clubs ON clubs.id = scored.club_id
            GROUP BY clubs.id, clubs.name
            ORDER BY goals DESC, clubs.name
            SQL, ['finished', 'finished']);

        return response()->json(['data' => $rows]);
    }
}
