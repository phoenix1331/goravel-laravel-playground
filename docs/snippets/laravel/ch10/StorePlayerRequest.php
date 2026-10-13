<?php

public function rules(): array
{
    return [
        'name' => ['required', 'string', 'max:100'],
        'position' => ['required', Rule::enum(Position::class)],
        'squad_number' => [
            'required', 'integer', 'min:1', 'max:99',
            // unique within this club only
            Rule::unique('players')->where('club_id', $this->route('club')->id),
        ],
    ];
}
