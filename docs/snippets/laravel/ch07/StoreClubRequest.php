<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class StoreClubRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true; // chapter 20 adds the real check
    }

    // Laravel's TrimStrings middleware already trims every input globally
    public function rules(): array
    {
        return [
            'name' => ['required', 'string', 'max:100', Rule::unique('clubs', 'name')],
            'ground' => ['required', 'string', 'max:100'],
            'founded_year' => ['nullable', 'integer', 'min:1857', 'max:'.now()->year],
        ];
    }

    public function messages(): array
    {
        return [
            'name.unique' => 'A club with this name already exists.',
            'founded_year.min' => 'The oldest football club was founded in 1857.',
        ];
    }
}
