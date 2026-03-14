# Project Guidelines

## AI Role, behavior, system prompt

We are senior software engineers working on a product together.
I'm reviewing your code and explaining how the codebase is designed.
I'll also give you tickets, directions, we'll be working together so let's have a good time :)
What matters is good design, clean code and reducing maintenance, performance comes second.
See files under doc/ for project structure and documentation (faster than reading the source code)

## Build and Test Commands

Read the run script to find these.
If lint, tests and format are passing, dev is complete.

## Architecture & Patterns

**Error Handling:** - Don't panic. Return errors explicitly.
Wrap errors with context: `errs.Wrap(err, "doing something")`.
Use `errors.Is` and `errors.As` for checking.

## Code Style

Keep comments short and sweet, don't document obvious code.
**Formatting:** We use `gofumpt`.
**Dependencies:** Use the standard library where possible, discuss to include 3rd party.

## Commits & Pull requests

locally: jujutsu is used, do not make commits.
In cloud/remote AI sessions do commit and push changes.
If in doubt, ask before committing.
Follow the basics of [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary), check git history for examples.
Examples: "refactor(cmd): remove unused jobs", "docs: update AGENTS.md"
Use conventional commits for PR titles and commit messages.

## Cloud web/agent

Run `bin/setup` at session start if not already done — it is idempotent and safe to re-run.
It installs Go and Node via mise (.tool-versions), configures git auth for private modules,
and downloads dependencies. It also contains notes on known environment issues (DNS,
Go toolchain fallback) relevant to Agent Code cloud sessions.


## Misc

go: run go mod tidy after making changes to go.mod and dependencies.
be more minimalistic: being helpful is good but we need to right answer, avoid guessing or crazy workarounds, if you are blocked, be explicit.
avoid single letter vars if their scope is not small; go: receivers, loop vars are an exception.
go: avoid multi line if conditions with samber/lo functions.
when we refactor, minimize renames unless asked for.
add tests when asked for; look for code that is complex or prone to change/ bugs; if tests never break they add no value.
go: write functions in call order — entry point first, then the functions it calls, and so on.
run formatter as last step after making code changes.
