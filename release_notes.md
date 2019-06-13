- Stopped building other services for integration tests.
- Added documentation in the form of updating the README and putting comments       
  in the yaml file.
- Changed how we delete: now we use the batchDeleter from the `codex` repo.  It 
  queries the database for ids of records that have passed their deathdate, 
  queues batches of expired records, then deletes them at a configurable rate.
