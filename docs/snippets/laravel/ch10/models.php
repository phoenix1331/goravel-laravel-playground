<?php

// app/Models/Club.php
public function players(): HasMany
{
    return $this->hasMany(Player::class);
}

// app/Models/Player.php
public function club(): BelongsTo
{
    return $this->belongsTo(Club::class);
}

// app/Models/Fixture.php: two foreign keys to the same table
public function homeClub(): BelongsTo
{
    return $this->belongsTo(Club::class, 'home_club_id');
}

public function awayClub(): BelongsTo
{
    return $this->belongsTo(Club::class, 'away_club_id');
}
