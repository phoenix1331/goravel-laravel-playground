<?php

public function test_records_results_from_the_feed(): void
{
    $this->artisan('sync:fixtures', ['--file' => base_path('tests/fixtures/fixture_results.json')])
        ->assertSuccessful();

    $this->assertTrue($this->fixture->fresh()->is_finished);
}

public function test_dry_run_changes_nothing(): void
{
    $this->artisan('sync:fixtures', ['--file' => base_path('tests/fixtures/fixture_results.json'), '--dry-run' => true])
        ->assertSuccessful();

    $this->assertSame('scheduled', $this->fixture->fresh()->status->value);
}

public function test_missing_file_fails(): void
{
    $this->artisan('sync:fixtures', ['--file' => 'missing.json'])->assertFailed();
}
