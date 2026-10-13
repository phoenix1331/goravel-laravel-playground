<?php

namespace App\Http\Controllers;

use App\Models\Club;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Response;

class ClubController extends Controller
{
    public function index(): JsonResponse
    {
        return response()->json(['data' => Club::orderBy('name')->get()]);
    }

    // Route model binding: Laravel loads the club, or returns 404 for you
    public function show(Club $club): JsonResponse
    {
        return response()->json(['data' => $club]);
    }

    public function store(Request $request): JsonResponse
    {
        // Validation arrives in chapter 07
        $club = Club::create($request->only('name', 'ground', 'founded_year'));

        return response()->json(['data' => $club], 201);
    }

    public function update(Request $request, Club $club): JsonResponse
    {
        $club->update($request->only('name', 'ground', 'founded_year'));

        return response()->json(['data' => $club]);
    }

    public function destroy(Club $club): Response
    {
        $club->delete();

        return response()->noContent();
    }
}
