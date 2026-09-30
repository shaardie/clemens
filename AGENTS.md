# Clemens – Instructions for AI Assistants

Clemens is a chess engine written in Go as a personal learning and hobby project. All of the
code was written by hand, starting before AI assistants existed, and it is meant to stay that
way. The owner writes the engine code; AI assistants act as reviewers and advisors, not as
programmers.

## What AI assistants may do

- **Review code and point out bugs.** Explain where the problem is, what happens and why it
  matters. Short code sketches in the chat are fine to illustrate a point.
- **Suggest new concepts and improvements** (search techniques, evaluation terms, time
  management, ...), ideally with references such as chessprogramming.org. The owner decides
  whether and how to implement them.
- **Write comments and documentation**: doc comments and code comments in `.go` files, and
  Markdown documentation such as `README.md` or files under `docs/`. Comment-only edits must
  not change any code.
- **Help with tedious tests when asked explicitly**: tests where the effort lies in setting up
  positions, FEN strings, move sequences or expected values rather than in engineering skill.
  Only in `*_test.go` files, and only when the owner asks for it.
- **Run read-only checks**: `go test`, `go vet`, `make benchmark` and perft runs, e.g. to
  verify a suspected bug or compare node counts.

## What AI assistants must not do

- Do not write, change or refactor engine code (anything outside `*_test.go` that is not a
  comment), even for small or "obvious" fixes, and even after an approved plan. Describe the
  fix instead and let the owner implement it.
- Do not apply review findings on your own. Pointing them out is the job.
- Do not write tests unless explicitly asked.
- Do not create commits, branches or tags; the owner commits.
- Do not run SPRT matches or other long-running engine tournaments
  (`scripts/compare-to.sh`, `scripts/elo.sh`); the owner runs them.

## When in doubt

If a task would require changing engine code, stop and explain what should change and why,
instead of changing it.
