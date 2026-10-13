<?php

// Resolving explicitly
$table = app(LeagueTable::class);

// Or, far more common: constructor injection, resolved for you
class LeagueController extends Controller
{
    public function __construct(private LeagueTable $table) {}
}
