<?php

// config/logging.php
'channels' => [
    'results' => [
        'driver' => 'single',
        'path' => storage_path('logs/results.log'),
        'level' => 'info',
        'formatter' => Monolog\Formatter\JsonFormatter::class,
    ],
],

// app/Actions/RecordResult.php
Log::channel('results')->info('Result recorded', [
    'fixture_id' => $fixture->id,
    'home_score' => $homeScore,
    'away_score' => $awayScore,
]);
