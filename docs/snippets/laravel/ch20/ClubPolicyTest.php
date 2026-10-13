<?php

public function test_club_policy_update(): void
{
    $club = new Club(['id' => 7]);
    $policy = new ClubPolicy();

    $this->assertTrue($policy->update(new User(['club_id' => 7]), $club)->allowed());
    $this->assertFalse($policy->update(new User(['club_id' => 8]), $club)->allowed());
    $this->assertFalse($policy->update(new User(), $club)->allowed());
}
