<?php

// app/Models/Fixture.php - a mutator that always stores kick-off in UTC
protected function kickOffAt(): Attribute
{
    return Attribute::set(fn ($value) => Carbon::parse($value)->utc());
}
