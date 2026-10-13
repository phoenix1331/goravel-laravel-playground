<?php

// app/Listeners/NotifyClubAdmins.php
public function handle(FixtureResultRecorded $event): void
{
    SendResultNotifications::dispatch($event->fixture);
}
