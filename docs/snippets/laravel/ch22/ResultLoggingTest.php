<?php

public function test_recording_a_result_writes_a_log_entry(): void
{
    Log::shouldReceive('channel')->with('results')->andReturnSelf();
    Log::shouldReceive('info')->once()->with('Result recorded', [
        'fixture_id' => $this->fixture->id,
        'home_score' => 4,
        'away_score' => 2,
    ]);

    app(RecordResult::class)($this->fixture, 4, 2);
}
