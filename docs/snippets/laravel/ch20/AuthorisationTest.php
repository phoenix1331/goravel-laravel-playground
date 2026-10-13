<?php

public function test_admins_cannot_edit_another_club(): void
{
    Sanctum::actingAs($this->awayAdmin);

    $this->putJson("/api/clubs/{$this->home->id}", ['name' => 'Hijacked', 'ground' => 'Park'])
        ->assertForbidden()
        ->assertExactJson(['message' => 'You can only manage your own club.']);
}

public function test_only_the_home_admin_records_the_result(): void
{
    Sanctum::actingAs($this->awayAdmin);
    $this->postJson("/api/fixtures/{$this->fixture->id}/result", ['home_score' => 0, 'away_score' => 5])
        ->assertForbidden();

    Sanctum::actingAs($this->homeAdmin);
    $this->postJson("/api/fixtures/{$this->fixture->id}/result", ['home_score' => 0, 'away_score' => 5])
        ->assertOk();
}
