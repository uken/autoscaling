# Auto Scaling tools

## How does it work?

We've split deploys into 3 phases:
- Build
- Housekeeping
- Release

For a rails app, it could mean:
- Build = `bundle install` `assets:precompile`
- Housekeeping = `db:migrate` `db:seed`
- Release = notify workers to load new app

These steps are streamlined via `sheepit`.

## Base Images

At the moment we have 2:
- uken/precise
Very similar to heroku's cedar stack. It does include a custom docker entrypoint required by `sheepit`.
- uken/frontend_19
Basic nginx + ruby 1.9.3 + nodejs

## Requirements

- Docker
- Consul
