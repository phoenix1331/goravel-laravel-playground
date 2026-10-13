<?php

namespace App\Listeners;

use App\Events\FixtureResultRecorded;
use App\Models\ResultNotification;
use App\Models\User;

class NotifyClubAdmins
{
    // Type-hinting the event is enough: Laravel discovers this listener
    public function handle(FixtureResultRecorded $event): void
    {
        $fixture = $event->fixture->load('homeClub', 'awayClub');
        $message = "Result: {$fixture->homeClub->name} {$fixture->home_score}-{$fixture->away_score} {$fixture->awayClub->name}";

        User::whereIn('club_id', [$fixture->home_club_id, $fixture->away_club_id])
            ->each(fn (User $admin) => ResultNotification::create([
                'user_id' => $admin->id,
                'fixture_id' => $fixture->id,
                'message' => $message,
            ]));
    }
}
