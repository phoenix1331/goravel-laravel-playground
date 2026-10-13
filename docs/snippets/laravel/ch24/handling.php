<?php

// Option 1: catch it where it matters
try {
    app(RecordResult::class)($fixture, $home, $away);
} catch (AlreadyRecorded $e) {
    return response()->json(['message' => 'This fixture already has a result.'], 409);
}

// Option 2: let it bubble up and render it centrally, in bootstrap/app.php
->withExceptions(function (Exceptions $exceptions) {
    $exceptions->render(fn (AlreadyRecorded $e) => response()->json(
        ['message' => 'This fixture already has a result.'], 409,
    ));
})
