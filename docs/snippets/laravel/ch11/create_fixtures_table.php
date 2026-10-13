<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('fixtures', function (Blueprint $table) {
            $table->id();
            $table->foreignId('home_club_id')->constrained('clubs')->cascadeOnDelete();
            $table->foreignId('away_club_id')->constrained('clubs')->cascadeOnDelete();
            $table->timestampTz('kick_off_at');
            $table->string('venue');
            $table->string('status', 20)->default('scheduled');
            $table->unsignedTinyInteger('home_score')->nullable();
            $table->unsignedTinyInteger('away_score')->nullable();
            $table->timestampsTz();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('fixtures');
    }
};
