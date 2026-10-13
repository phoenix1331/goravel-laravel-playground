<?php

namespace Database\Factories;

use Illuminate\Database\Eloquent\Factories\Factory;

class ClubFactory extends Factory
{
    public function definition(): array
    {
        $town = fake()->randomElement(['Ashford', 'Barnsley', 'Crewe', /* ... */]);

        return [
            'name' => $town.' '.fake()->randomElement(['Athletic', 'Rovers', 'United', 'Town', 'Wanderers']),
            'ground' => $town.' '.fake()->randomElement(['Park', 'Rec', 'Lane', 'Road']),
            'founded_year' => fake()->numberBetween(1870, 2020),
        ];
    }
}

// app/Models/Club.php
class Club extends Model
{
    use HasFactory;
}
