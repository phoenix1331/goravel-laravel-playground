<?php

// ClubController::update - throws AuthorizationException (403) when denied
public function update(UpdateClubRequest $request, Club $club): JsonResponse
{
    Gate::authorize('update', $club);

    // ...
}

// Or on the route: Route::put('/clubs/{club}', ...)->can('update', 'club');
