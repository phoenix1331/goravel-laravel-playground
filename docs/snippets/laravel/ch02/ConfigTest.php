<?php

namespace Tests\Unit;

use Tests\TestCase;

class ConfigTest extends TestCase
{
    public function test_typed_lookups(): void
    {
        // phpunit.xml / .env.testing overrides the default
        $this->assertSame('Test League', config('club.league_name'));

        // config() returns mixed: the (int) cast lives in config/club.php
        $this->assertSame(3, config('club.points.win'));
        $this->assertSame(1, config('club.points.draw'));
        $this->assertSame(0, config('club.points.loss'));

        // Missing keys fall back to the default passed in
        $this->assertSame(42, config('club.does_not_exist', 42));
    }
}
