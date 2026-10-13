<?php

namespace App\Services;

use Illuminate\Support\Facades\Http;

class HttpFeed implements FixtureFeed
{
    public function results(): array
    {
        return Http::baseUrl(config('services.fixtures.url'))
            ->timeout(10)
            ->acceptJson()
            ->get('/results')
            ->throw()   // 4xx and 5xx become a RequestException
            ->json();
    }
}
