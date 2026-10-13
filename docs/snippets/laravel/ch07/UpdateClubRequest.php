<?php

// UpdateClubRequest::rules(): ignore the club being updated
'name' => [
    'required', 'string', 'max:100',
    Rule::unique('clubs', 'name')->ignore($this->route('club')),
],

// or as a string: unique:table,column,ignoredValue,idColumn
'name' => 'required|string|max:100|unique:clubs,name,'.$this->route('club')->id.',id',
