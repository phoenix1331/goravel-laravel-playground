<?php

// routes/api.php
Route::post('/auth/login', [AuthController::class, 'login']);
Route::get('/clubs', [ClubController::class, 'index']);

Route::middleware('auth:sanctum')->group(function () {
    Route::get('/auth/me', [AuthController::class, 'me']);
    Route::post('/auth/logout', [AuthController::class, 'logout']);
    Route::post('/clubs', [ClubController::class, 'store']);
    // ...
});
