<?php

namespace App\Policies;

use App\Models\Club;
use App\Models\User;
use Illuminate\Auth\Access\Response;

// Found automatically for the Club model: no registration needed
class ClubPolicy
{
    public function create(User $user): Response
    {
        return $user->club_id === null
            ? Response::allow()
            : Response::deny('You already run a club.');
    }

    // Laravel passes the logged-in user as the first argument
    public function update(User $user, Club $club): Response
    {
        return $user->club_id === $club->id
            ? Response::allow()
            : Response::deny('You can only manage your own club.');
    }
}
