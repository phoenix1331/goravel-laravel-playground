<?php

namespace App\Models;

use Illuminate\Foundation\Auth\User as Authenticatable;

class User extends Authenticatable
{
    protected $fillable = ['club_id', 'name', 'email', 'password'];

    // Left out of toArray() / toJson()
    protected $hidden = ['password'];
}
