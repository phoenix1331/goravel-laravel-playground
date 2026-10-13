<?php

namespace Database\Seeders;

use App\Models\Club;
use App\Models\User;
use Illuminate\Database\Seeder;

class DatabaseSeeder extends Seeder
{
    public function run(): void
    {
        fake()->seed(2026);

        $clubs = Club::factory()->count(6)->create();

        $clubs->each(function (Club $club, int $i) {
            foreach (range(1, 16) as $number) {
                $club->players()->save(Player::factory()->make(['squad_number' => $number]));
            }

            User::factory()->for($club)->create(['email' => 'admin'.($i + 1).'@example.com']);
        });

        // ... a round of fixtures
    }
}
