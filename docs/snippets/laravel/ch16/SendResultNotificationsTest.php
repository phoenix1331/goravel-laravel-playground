<?php

public function test_job_handle_writes_notifications(): void
{
    (new SendResultNotifications($this->fixture))->handle();

    $this->assertDatabaseHas('result_notifications', [
        'user_id' => $this->admin->id,
        'message' => 'Result: Ashford 2-2 Barnsley',
    ]);
}

public function test_recording_a_result_queues_the_job(): void
{
    Queue::fake();

    $this->postJson("/api/fixtures/{$this->fixture->id}/result", ['home_score' => 1, 'away_score' => 0]);

    Queue::assertPushed(SendResultNotifications::class);
}
