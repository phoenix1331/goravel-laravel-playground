<?php

// ClubPlayerController::index - lazy eager load onto a model you already have
public function index(Club $club): JsonResponse
{
    $club->load(['players' => fn ($query) => $query->orderBy('squad_number')]);

    return response()->json(['data' => $club]);
}

// ClubPlayerController::store - create through the relation
public function store(StorePlayerRequest $request, Club $club): JsonResponse
{
    $player = $club->players()->create($request->validated());

    return response()->json(['data' => $player], 201);
}

// FixtureController::index - eager load while querying
public function index(): JsonResponse
{
    $fixtures = Fixture::with(['homeClub', 'awayClub'])->orderBy('kick_off_at')->get();

    return response()->json(['data' => $fixtures]);
}
