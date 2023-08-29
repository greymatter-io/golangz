# Golangz Changelog

## [vxxxx] -- 2023-08-09
- Adds Generic Either and Option types
- Fixes bug in sets.Union to remove duplicates

## [v0.1.21] -- 2023-03-01
- Removes go.mod because it is not necessary given go mod init and go mod tidy

## [v0.1.20] -- 2022-11-29
- Corrects sort to mutate array and not lie about it.

## [v0.1.16] -- 2022-07-01

### Changed
- Moves Set operations from arrays package to sets package. This is an API breaking change.
- Creates arrays equality function whic respects ordering

