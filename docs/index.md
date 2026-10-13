---
layout: home

hero:
  name: Goravel playground for Laravel devs
  text: Something to run, poke and break
  tagline: A small Goravel app for Laravel developers, with short notes on how each feature compares to Laravel.
  actions:
    - theme: brand
      text: Go for PHP devs
      link: /go-for-php-devs
    - theme: alt
      text: Chapter 01
      link: /chapters/01-installation
    - theme: alt
      text: Goravel docs
      link: https://www.goravel.dev/getting-started/installation.html

features:
  - title: A real app to play with
    details: A grassroots football club API with clubs, players, fixtures and a league table. Clone it, run it with Docker, change it and watch the tests.
  - title: Short notes, real code
    details: Each chapter shows the code in this repo next to a Laravel snippet. The Go is pulled from the source at build time and tested in CI.
  - title: Where Go is different
    details: Explicit errors, structs instead of classes, no magic methods, compile-time wiring. The notes call these out where they bite.
---

## What this is

An independent, unofficial playground for Laravel developers trying [Goravel](https://www.goravel.dev), which is built and maintained by the Goravel team and contributors. It is not affiliated with or endorsed by the Goravel project.

It is not a replacement for the official docs. For the full feature-by-feature comparison, read Goravel's [Compare with Laravel](https://www.goravel.dev/prologue/compare-with-laravel.html) page; for everything else, the [Goravel documentation](https://www.goravel.dev/getting-started/installation.html).

What this adds is a working app to experiment with. Goravel deliberately mirrors Laravel, and the notes here point out where the resemblance holds and where Go forces a different approach: a facade is not a static proxy, a model has no `__get`, an error is a value you return rather than an exception you throw.

## How to use it

1. Clone the repo and follow the quick start in the README: `make up`, `make fresh`, then open `http://localhost:8000/`.
2. Read [Go for PHP devs](/go-for-php-devs) if Go is new to you, and [Reading the comparisons](/reading-the-comparisons) for how each chapter is laid out.
3. Pick a chapter, find the code it points to, change something, and run `make test`.
