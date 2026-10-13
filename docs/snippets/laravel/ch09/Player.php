<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

enum Position: string
{
    case Goalkeeper = 'goalkeeper';
    case Defender = 'defender';
    case Midfielder = 'midfielder';
    case Forward = 'forward';
}

class Player extends Model
{
    // No properties declared: columns arrive at runtime via __get / __set
    protected $fillable = ['club_id', 'name', 'position', 'squad_number'];

    protected function casts(): array
    {
        return [
            'position' => Position::class,
            'squad_number' => 'integer',
        ];
    }
}
