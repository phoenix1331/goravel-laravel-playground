<?php

// routes/api.php - Laravel applies the "api" prefix in bootstrap/app.php

use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Route;

// A closure
Route::get('/league', fn () => ['name' => config('club.league_name')]);

// A named handler: Laravel wants a controller class (or an invokable),
// not a plain function
Route::get('/league/points', LeaguePointsController::class);
