<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\DB;
use Throwable;

class HealthController extends Controller
{
    public function show(): JsonResponse
    {
        try {
            DB::connection()->getPdo();
        } catch (Throwable) {
            return response()->json(['status' => 'error', 'database' => 'down'], 503);
        }

        return response()->json(['status' => 'ok', 'database' => 'up']);
    }
}
