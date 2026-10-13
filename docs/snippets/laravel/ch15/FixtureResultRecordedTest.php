<?php

public function test_recording_a_result_notifies_both_clubs_admins(): void
{
    $this->postJson("/api/fixtures/{$this->fixture->id}/result", ['home_score' => 3, 'away_score' => 1])
        ->assertOk();

    foreach (['Ashford', 'Barnsley'] as $club) {
        $this->assertDatabaseHas('result_notifications', [
            'user_id' => $this->admins[$club]->id,
            'message' => 'Result: Ashford 3-1 Barnsley',
        ]);
    }
    $this->assertDatabaseMissing('result_notifications', ['user_id' => $this->admins['Crewe']->id]);
}

// Or check only that the event fired, without running listeners
public function test_recording_a_result_fires_the_event(): void
{
    Event::fake();

    $this->postJson("/api/fixtures/{$this->fixture->id}/result", ['home_score' => 3, 'away_score' => 1]);

    Event::assertDispatched(FixtureResultRecorded::class);
}
