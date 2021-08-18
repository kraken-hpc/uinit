# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2021-08-18
### Changed
- Refactor into a reusable package with a runnable Script object
### Added
- Script module now allows execution of subscripts

## [0.1.3] - 2021-03-23
### Fixed
- Fix a critical bug that caused foreground apps that need PTYs to fail

## [0.1.2] - 2021-03-22
### Added
- Multiplex logs to stdio & log file
- Add simple loop capabilities
- Support for specifying Env vars in command module

### Fixed
- All process output recorded in logs
- Processes' output is prefixed in logs

## [0.1.1] - 2021-03-19
### Chnaged
- Migrate to `kraken-hpc` github org

## [0.1.0] - 2021-01-05
### Added
- Initial versioned release
