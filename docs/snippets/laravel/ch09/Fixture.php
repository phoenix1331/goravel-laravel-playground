<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Casts\Attribute;
use Illuminate\Database\Eloquent\Model;

class Fixture extends Model
{
    protected $fillable = [
        'home_club_id', 'away_club_id', 'kick_off_at', 'venue', 'status', 'home_score', 'away_score',
    ];

    protected function casts(): array
    {
        return [
            'kick_off_at' => 'datetime',
            'status' => FixtureStatus::class,
            'home_score' => 'integer',
            'away_score' => 'integer',
        ];
    }

    // $fixture->is_finished
    protected function isFinished(): Attribute
    {
        return Attribute::get(fn () => $this->status === FixtureStatus::Finished
            && $this->home_score !== null
            && $this->away_score !== null);
    }
}
