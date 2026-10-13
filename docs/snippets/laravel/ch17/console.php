<?php

// routes/console.php

use App\Models\ResultNotification;
use Illuminate\Support\Facades\Schedule;

Schedule::command('sync:fixtures')
    ->everyFifteenMinutes()
    ->saturdays()
    ->withoutOverlapping();

Schedule::call(fn () => ResultNotification::where('created_at', '<', now()->subDays(30))->delete())
    ->dailyAt('03:00')
    ->name('prune-old-notifications');
