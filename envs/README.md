# Sweet Tea Environment Variables

This file outlines the environment variables used in Sweet Tea.

* `API_CLUSTER_ZONES`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.ApiClusterZones`
	- Description: TODO
	- Ex: `us-west-1a`

* `API_VERSION`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.ApiVersion`
	- Description: Version of the SweetTea API to use. All API routes will be prefixed with this value
	- Ex: `v1`
	 
* `AWS_ACCESS_KEY_ID`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.AwsAccessKeyId`
	- Description: TODO
	- Ex: `AKABC123DEF456GHI789`

* `AWS_REGION_NAME`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.AwsRegionName`
	- Description: TODO 
	- Ex: `us-west-1`

* `AWS_SECRET_ACCESS_KEY`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.AwsSecretAccessKey`
	- Description: TODO
	- Ex: `ay+aBc123dEf/456gHi7/89JkLmNoPqRsTuVwXyZ`
	 
* `BUILD_CLUSTER_NAME`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.BuildClusterName`
	- Description: TODO
	- Ex: `build-local`

* `BUILD_CLUSTER_STATE`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.BuildClusterState`
	- Description: TODO
	- Ex: `s3://my-bucket`

* `BUILD_CLUSTER_ZONES`:
	- Required to start app: No
	- Type: `string`
	- In-app ref: None
	- Description: TODO
	- Ex: `us-west-1a`
	 
* `CLOUD_PROVIDER`:
	- Required to start app: Yes
	- Type: `string`
	- Accepted values: `aws`
	- In-app ref: `app.Config.CloudProvider`
	- Description: TODO
	- Ex: `aws`

* `CLUSTER_IMAGE`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.ClusterImage`
	- Description: TODO 
	- Ex: `099720109477/ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20171026.1`

* `CORE_CLUSTER_NAME`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.CoreClusterName`
	- Description: TODO
	- Ex: `core-local`
	 
* `CORE_CLUSTER_STATE`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.CoreClusterState`
	- Description: TODO
	- Ex: `s3://my-bucket`

* `CORE_CLUSTER_ZONES`:
	- Required to start app: No
	- Type: `string`
	- In-app ref: None
	- Description: TODO
	- Ex: `us-west-1a`

* `DATABASE_URL`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.DatabaseUrl`
	- Description: PostgreSQL database URL
	- Ex: `postgres://username:password@host:port/username`

* `DEBUG`:
	- Required to start app: Yes
	- Type: `bool`
	- In-app ref: `app.Config.Debug`
	- Description: Whether to log database & API calls more verbosely for debugging purposes
	- Ex: `true`
	 
* `DOMAIN`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.Domain`
	- Description: TODO
	- Ex: `mysite.com`

* `ENV`:
	- Required to start app: Yes
	- Type: `string`
	- Accepted values: `test|local|dev|staging|prod`
	- In-app ref: `app.Config.Env`
	- Description: The environment tier currently running
	- Ex: `local`
	 
* `HOSTED_ZONE_ID`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.HostedZoneId`
	- Description: TODO
	- Ex: `ZABC123DEF4567`

* `IMAGE_OWNER`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.ImageOwner`
	- Description: TODO
	- Ex: `my-app-local`

* `IMAGE_OWNER_PW`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.ImageOwnerPw`
	- Description: TODO
	- Ex: `abc123`
	 
* `JOB_QUEUE_NSP`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.JobQueueNsp`
	- Description: Redis namespace used for enqueueing jobs for the worker
	- Ex: `my_job_queue`

* `K8S_VERSION`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.K8sVersion`
	- Description: TODO
	- Ex: `1.7.10`

* `KUBECONFIG`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.KubeConfig`
	- Description: TODO
	- Ex: `/root/.kubeconfig`
	 
* `MASTER_SIZE`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.MasterSize`
	- Description: TODO
	- Ex: `t2.small`

* `MIGRATE_IMAGE_NAME`:
	- Required to start app: No
	- Type: `string`
	- In-app ref: None
	- Description: TODO
	- Ex: `my-app-migrate`

* `MIGRATE_RESTART_POLICY`:
	- Required to start app: No
	- Type: `string`
	- Accepted values: `Always|OnFailure|Never`
	- In-app ref: None
	- Description: TODO
	 
* `MIGRATE_REPLICAS`:
	- Required to start app: No
	- Type: `int`
	- In-app ref: None
	- Description: TODO
	- Ex: `1`

* `NODE_COUNT`:
	- Required to start app: Yes
	- Type: `int`
	- In-app ref: `app.Config.NodeCount`
	- Description: TODO
	- Ex: `2`

* `NODE_SIZE`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.NodeSize`
	- Description: TODO
	- Ex: `t2.small`
	 
* `REDIS_POOL_MAX_ACTIVE`:
	- Required to start app: Yes
	- Type: `int`
	- In-app ref: `app.Config.RedisPoolMaxActive`
	- Description: Max number of connections allocated by the Redis pool at any given time (0 = unlimited)
	- Ex: `5`

* `REDIS_POOL_MAX_IDLE`:
	- Required to start app: Yes
	- Type: `int`
	- In-app ref: `app.Config.RedisPoolMaxIdle`
	- Description: Max number of idle connections in the Redis pool
	- Ex: `5`

* `REDIS_POOL_WAIT`:
	- Required to start app: Yes
	- Type: `bool`
	- In-app ref: `app.Config.RedisPoolWait`
	- Description: Whether to wait for new connections to become available if the Redis pool is at its max connection limit
	- Ex: `true`
	 
* `REDIS_URL`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.RedisUrl`
	- Description: Redis URL
	- Ex: `localhost:6379`

* `REST_API_TOKEN`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.RestApiToken`
	- Description: TODO
	- Ex: `abc-123-def-456`

* `SERVER_IMAGE_NAME`:
	- Required to start app: No
	- Type: `string`
	- In-app ref: None
	- Description: TODO
	- Ex: `my-app-server`
	 
* `SERVER_PORT`:
	- Required to start app: Yes
	- Type: `int`
	- In-app ref: `app.Config.ServerPort`
	- Description: Port the server runs on
	- Ex: `80`

* `SERVER_RESTART_POLICY`:
	- Required to start app: No
	- Type: `string`
	- Accepted values: `Always|OnFailure|Never`
	- In-app ref: None
	- Description: TODO

* `SERVER_REPLICAS`:
	- Required to start app: No
	- Type: `int`
	- In-app ref: None
	- Description: TODO
	- Ex: `3`
	 
* `TRAIN_CLUSTER_NAME`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.TrainClusterName`
	- Description: TODO
	- Ex: `train-local`

* `TRAIN_CLUSTER_STATE`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.TrainClusterState`
	- Description: TODO
	- Ex: `s3://my-bucket`

* `TRAIN_CLUSTER_ZONES`:
	- Required to start app: No
	- Type: `string`
	- In-app ref: None
	- Description: TODO
	- Ex: `us-west-1a`

* `WILDCARD_SSL_CERT_ARN`:
	- Required to start app: Yes
	- Type: `string`
	- In-app ref: `app.Config.WildcardSslCertArn`
	- Description: TODO
	- Ex: `arn:aws:acm:us-west-1:123456789123:certificate/abc123de-f456-ghi7-89jk-lmnopqrstuvw`

* `WORKER_COUNT`:
	- Required to start app: Yes
	- Type: `int`
	- In-app ref: `app.Config.WorkerCount`
	- Description: Number of workers to allocate for the worker pool
	- Ex: `10`

* `WORKER_IMAGE_NAME`:
	- Required to start app: No
	- Type: `string`
	- In-app ref: None
	- Description: TODO
	- Ex: `my-app-worker`

* `WORKER_RESTART_POLICY`:
	- Required to start app: No
	- Type: `string`
	- Accepted values: `Always|OnFailure|Never`
	- In-app ref: None
	- Description: TODO
	 
* `WORKER_REPLICAS`:
	- Required to start app: No
	- Type: `int`
	- In-app ref: None
	- Description: TODO
	- Ex: `2`