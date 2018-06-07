# Sweet Tea Environment Variables

This file outlines the environment variables used in Sweet Tea, the roles the play, their defaults,
and if they are required.

* `ST_DATABASEURL`:
	- Required: Yes
	- Type: `string`
	- Default: None
	- In-app accessibility: `app.Config.DatabaseUrl`
	- Description: PostgreSQL database URL

* `ST_REDISURL`:
	- Required: Yes
	- Type: `string`
	- Default: None
	- In-app accessibility: `app.Config.RedisUrl`
	- Description: Redis URL

* `ST_APIVERSION`:
	- Required: No
	- Type: `string`
	- Default: `v1`
	- In-app accessibility: `app.Config.ApiVersion`
	- Description: Version of the SweetTea API to use. All API routes will be prefixed with this value

* `ST_DEBUG`:
	- Required: No
	- Type: `bool`
	- Values: `true|false`
	- Default: `false`
	- In-app accessibility: `app.Config.Debug`
	- Description: Whether to log database & API calls more verbosely for debugging purposes

* `ST_ENV`:
	- Required: No
	- Type: `string`
	- Values: `test|local|dev|staging|prod`
	- Default: `local`
	- In-app accessibility: `app.Config.Env`
	- Description: The environment tier currently running

* `ST_JOBQUEUENSP`:
	- Required: No
	- Type: `string`
	- Default: `st_job_queue`
	- In-app accessibility: `app.Config.JobQueueNsp`
	- Description: Redis namespace used for enqueueing jobs for the worker

* `ST_PORT`:
	- Required: No
	- Type: `int`
	- Default: `5000`
	- In-app accessibility: `app.Config.Port`
	- Description: Port the server runs on

* `ST_REDISPOOLMAXACTIVE`:
	- Required: No
	- Type: `int`
	- Default: `5`
	- In-app accessibility: `app.Config.RedisPoolMaxActive`
	- Description: Max number of connections allocated by the Redis pool at any given time (0 = unlimited)

* `ST_REDISPOOLMAXIDLE`:
	- Required: No
	- Type: `int`
	- Default: `5`
	- In-app accessibility: `app.Config.RedisPoolMaxIdle`
	- Description: Max number of idle connections in the Redis pool

* `ST_REDISPOOLWAIT`:
	- Required: No
	- Type: `bool`
	- Default: `true`
	- In-app accessibility: `app.Config.RedisPoolWait`
	- Description: Whether to wait for new connections to become available if the Redis pool is at its max connection limit

* `ST_WORKERCOUNT`:
	- Required: No
	- Type: `int`
	- Default: `10`
	- In-app accessibility: `app.Config.WorkerCount`
	- Description: Number of workers to allocate for the worker pool
