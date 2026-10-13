<?php

public function test_forget_league_table_forgets_the_cached_table(): void
{
    Cache::shouldReceive('forget')->once()->with('league_table')->andReturnTrue();

    (new ForgetLeagueTable())->handle(new FixtureResultRecorded($this->fixture));
}
