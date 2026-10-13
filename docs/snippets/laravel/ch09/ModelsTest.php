<?php

namespace Tests\Feature;

use App\Models\Club;
use App\Models\Fixture;
use App\Models\Player;
use App\Models\Position;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class ModelsTest extends TestCase
{
    use RefreshDatabase;

    public function test_create_fills_id_and_timestamps(): void
    {
        $club = Club::create(['name' => 'Ashford Athletic', 'ground' => 'The Rec']);

        $this->assertNotNull($club->id);
        $this->assertNotNull($club->created_at);
    }

    public function test_typed_fields_round_trip(): void
    {
        $club = Club::create(['name' => 'Ashford Athletic', 'ground' => 'The Rec']);
        Player::create(['club_id' => $club->id, 'name' => 'Sam Keeper', 'position' => Position::Goalkeeper, 'squad_number' => 1]);

        $found = Player::where('squad_number', 1)->firstOrFail();

        $this->assertSame(Position::Goalkeeper, $found->position);
        $this->assertSame(1, $found->squad_number);
    }

    public function test_password_is_never_serialised(): void
    {
        $user = new User(['name' => 'Alex Admin', 'email' => 'alex@example.com', 'password' => 'hashed']);

        $this->assertStringNotContainsString('password', $user->toJson());
    }
}
