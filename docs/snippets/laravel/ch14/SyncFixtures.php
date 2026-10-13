<?php

namespace App\Console\Commands;

use App\Actions\RecordResult;
use App\Services\FileFeed;
use Illuminate\Console\Command;

class SyncFixtures extends Command
{
    // Name, arguments and options in one string
    protected $signature = 'sync:fixtures {--file= : Read results from this JSON file} {--dry-run}';

    protected $description = "Record results from the league's fixture feed";

    public function handle(RecordResult $record): int
    {
        $path = $this->option('file') ?: config('club.feed.path');

        try {
            $results = (new FileFeed($path))->results();
        } catch (\RuntimeException $e) {
            $this->error($e->getMessage());

            return self::FAILURE;
        }

        $rows = [];
        foreach ($results as $result) {
            // ... find the open fixture, record unless --dry-run
        }

        $this->table(['Home', 'Away', 'Score', 'Outcome'], $rows);

        return self::SUCCESS;
    }
}
