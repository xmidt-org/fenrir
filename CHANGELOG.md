# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.5.3]




## [v0.5.2]
- updated urls and imports



## [v0.5.1]




## [v0.5.0]
- bumped codex common to change queue to set



## [v0.4.0]
- bumped Codex for batchDeleter to use Nanosecond Unix time



## [v0.3.2]
- bumped codex



## [v0.3.1]
- bumped codex



## [v0.3.0]
- Stopped building other services for integration tests.
- Added documentation in the form of updating the README and putting comments       
  in the yaml file.
- Changed how we delete: now we use the batchDeleter from the `codex` repo.  It 
  queries the database for ids of records that have passed their deathdate, 
  queues batches of expired records, then deletes them at a configurable rate.



## [v0.2.0]
- Added Prune limit
- Leverage mutliple connections
- bumped codex to v0.3.2
- bumped webpa-common to v1.0.0



## [v0.1.1]
- added metrics
- added pprof



## [v0.1.0]
- Inital Release

[Unreleased]: https://github.com/xmidt-org/fenrir/compare/v0.5.3...HEAD
[v0.5.3]: https://github.com/xmidt-org/fenrir/compare/v0.5.2...v0.5.3
[v0.5.2]: https://github.com/xmidt-org/fenrir/compare/v0.5.1...v0.5.2
[v0.5.1]: https://github.com/xmidt-org/fenrir/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/xmidt-org/fenrir/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/xmidt-org/fenrir/compare/v0.3.2...v0.4.0
[v0.3.2]: https://github.com/xmidt-org/fenrir/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/xmidt-org/fenrir/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/xmidt-org/fenrir/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/xmidt-org/fenrir/compare/v0.1.1...v0.2.0
[v0.1.1]: https://github.com/xmidt-org/fenrir/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/xmidt-org/fenrir/compare/v0.0.0...v0.1.0
