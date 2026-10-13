<?php

// config/club.php

return [
    // Name shown on the league table
    'league_name' => env('LEAGUE_NAME', 'Sunday League Division One'),

    // League points awarded per result
    'points' => [
        'win' => (int) env('LEAGUE_POINTS_WIN', 3),
        'draw' => (int) env('LEAGUE_POINTS_DRAW', 1),
        'loss' => (int) env('LEAGUE_POINTS_LOSS', 0),
    ],

    // Largest squad a club may register
    'max_squad_size' => (int) env('CLUB_MAX_SQUAD_SIZE', 25),
];
