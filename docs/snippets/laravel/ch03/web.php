<?php

// routes/web.php
// Laravel 11+ also ships a built-in `health: '/up'` route in bootstrap/app.php.
// This explicit route mirrors the Goravel version, which checks the database.

use App\Http\Controllers\HealthController;
use Illuminate\Support\Facades\Route;

Route::get('/health', [HealthController::class, 'show']);
