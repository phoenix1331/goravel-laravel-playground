<?php

namespace App\Jobs;

use App\Models\Fixture;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;

class SendResultNotifications implements ShouldQueue
{
    use Queueable;

    // The model is serialised as its ID and re-fetched by the worker
    public function __construct(public Fixture $fixture) {}

    public function handle(): void
    {
        // ... write a notification for each admin of both clubs
    }
}
