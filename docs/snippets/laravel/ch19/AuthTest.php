<?php

public function test_login_returns_a_token(): void
{
    $this->postJson('/api/auth/login', ['email' => 'alex@example.com', 'password' => 'password'])
        ->assertOk()
        ->assertJsonPath('token_type', 'Bearer');
}

public function test_me_returns_the_tokens_user(): void
{
    Sanctum::actingAs($this->user);

    $this->getJson('/api/auth/me')
        ->assertOk()
        ->assertJsonPath('data.email', 'alex@example.com')
        ->assertJsonMissingPath('data.password');
}

public function test_writes_need_a_token(): void
{
    $this->postJson('/api/clubs', ['name' => 'Anon FC', 'ground' => 'Park'])->assertUnauthorized();
}
