<?php

// routes/api.php
Route::middleware(EnsureJsonRequest::class)->group(function () {
    Route::apiResource('clubs', ClubController::class);
});

// Or globally, in bootstrap/app.php
->withMiddleware(function (Middleware $middleware): void {
    $middleware->api(append: [EnsureJsonRequest::class]);
})
