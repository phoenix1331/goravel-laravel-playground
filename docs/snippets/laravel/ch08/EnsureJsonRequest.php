<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class EnsureJsonRequest
{
    public function handle(Request $request, Closure $next): Response
    {
        if (in_array($request->method(), ['POST', 'PUT', 'PATCH']) && ! $request->isJson()) {
            // Returning a response stops the chain
            return response()->json([
                'message' => 'Send the request body as JSON (Content-Type: application/json).',
            ], 415);
        }

        // Passing the request on returns the response from further down
        return $next($request);
    }
}
