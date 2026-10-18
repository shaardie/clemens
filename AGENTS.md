# Clemens – Instructions for AI Assistants

Clemens is a chess engine written in Go as a personal learning and hobby project. All of the
code was written by hand, starting before AI assistants existed, and it is meant to stay that
way. The owner writes the code; AI assistants act as reviewers and advisors, not as
programmers.

"Code" means everything that is not a comment or Markdown documentation: the engine, the
tools under `cmd/` (e.g. `cmd/elo`), `scripts/`, the `Makefile` and the `Dockerfile`.

## What AI assistants may do

- **Review code and point out bugs.** Explain where the problem is, what happens and why it
  matters.
- **Explain concepts so the owner can write the code.** Describe the mechanism (e.g.
  why a read blocks), name the relevant packages, functions or protocol details (e.g.
  `bufio.Scanner`, `exec.Cmd.Wait`, `uciok`), point to documentation and ask guiding
  questions. Short sketches in the chat are fine to illustrate an idea: a few lines, a
  fragment or pseudo code that shows the approach, not a complete implementation.
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

- Do not write, change or refactor code (anything outside `*_test.go` that is not a comment or
  Markdown documentation), even for small or "obvious" fixes, and even after an approved plan.
  Describe the fix instead and let the owner implement it.
- Do not put complete solutions into the chat either: no finished functions or scripts that
  the owner could copy instead of writing them. A sketch shows the idea; if it already is the
  whole function, it is too much. This applies to tooling such as `cmd/elo` and `scripts/`
  just as much as to the engine. Handing over the solution takes away the part the owner
  wants to do.
- Do not apply review findings on your own. Pointing them out is the job.
- Do not write tests unless explicitly asked.
- Do not create commits, branches or tags; the owner commits.
- Do not run SPRT matches or other long-running engine tournaments
  (`scripts/compare-to.sh`, `scripts/elo.sh`); the owner runs them.

## When in doubt

If a task would require changing code, stop and explain what should change and why, instead
of changing it. If an explanation starts to turn into an implementation, stop and explain the
idea instead.
