<?php

public function test_sync_runs_every_fifteen_minutes_on_saturdays(): void
{
    $event = collect(app(Schedule::class)->events())
        ->first(fn ($event) => str_contains($event->command, 'sync:fixtures'));

    $this->assertSame('*/15 * * * 6', $event->expression);
    $this->assertTrue($event->withoutOverlapping);
}
