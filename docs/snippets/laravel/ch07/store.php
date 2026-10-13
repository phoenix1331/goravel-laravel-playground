<?php

// ClubController: type-hint the request and validation has already happened.
// A failure never reaches this method: Laravel throws a ValidationException
// and the exception handler turns it into a 422.
public function store(StoreClubRequest $request): JsonResponse
{
    $club = Club::create($request->validated());

    return response()->json(['data' => $club], 201);
}
