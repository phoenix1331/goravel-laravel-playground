---
layout: home

hero:
  name: Goravel for Laravel devs
  text: Learn Go through the framework you already know
  tagline: Every Laravel concept, mapped to its Goravel equivalent, side by side, backed by a working Go app.
  actions:
    - theme: brand
      text: Go for PHP devs
      link: /go-for-php-devs
    - theme: alt
      text: Chapter 01
      link: /chapters/01-installation
    - theme: alt
      text: Cheat sheet
      link: /cheat-sheet

features:
  - title: Side by side
    details: Each chapter puts Laravel code next to the Goravel code that does the same job, then explains what is different and why.
  - title: Real, tested code
    details: The Go snippets are pulled from the repo's source at build time and run in CI against SQLite and MySQL. No pseudo-code.
  - title: Where Go forces a different approach
    details: Explicit errors, structs instead of classes, no magic methods, compile-time wiring. The chapters call these out where they bite.
---

## Why this exists

Goravel deliberately mirrors Laravel: artisan, facades, Eloquent-style ORM, migrations, queues. That makes it a gentle way into Go for a Laravel developer, but the resemblance hides real differences. A facade is not a static proxy. A model has no `__get`. An error is a value you return, not an exception you throw.

This site walks through 24 Laravel concepts one at a time. The code behind it is a small grassroots football club API: clubs, players, fixtures and club admins who can only manage their own club.

## How to use it

1. Read [Go for PHP devs](/go-for-php-devs) if Go is new to you.
2. See [Reading the comparisons](/reading-the-comparisons) for how each chapter is laid out.
3. Work through the chapters in order, or jump in using the [cheat sheet](/cheat-sheet).
