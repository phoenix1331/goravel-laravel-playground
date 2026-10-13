<?php

// routes/web.php - optional in Laravel: unknown URLs already return 404

Route::fallback(fn () => response()->json(['message' => 'Not Found'], 404));
