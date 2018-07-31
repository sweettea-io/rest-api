# SweetTea Database Schema

Models:

    User
    Session
    Project
    ProjectConfig
    Commit
    TrainJob
    Model
    ModelVersion
    Deploy
    ApiCluster
    EnvVar

Relationships:

    User|Session
        User --> has many --> Sessions
        Session --> belongs to --> User

    Project|Commit
        Project --> has many --> Commits
        Commit --> belongs to --> Project

    Project|ProjectConfig
        Project --> has one --> ProjectConfig
        ProjectConfig --> has one --> Project

    Project|Model
        Project --> has many --> Models
        Model --> belongs to --> Project

    Commit|TrainJob
        Commit --> has many --> TrainJobs
        TrainJob --> belongs to --> Commit

    Commit|Deploy
        Commit --> has many --> Deploys
        Deploy --> belongs to --> Commit

    Model|ModelVersion
        Model --> has many --> ModelVersions
        ModelVersion --> belongs to --> Model

    TrainJob|ModelVersion
        TrainJob --> has one --> ModelVersion
        ModelVersion --> has one --> TrainJob

    ModelVersion|Deploy
        ModelVersion --> has many --> Deploys
        Deploy --> belongs to --> ModelVersion

    ApiCluster|Deploy
        ApiCluster --> has many --> Deploys
        Deploy --> belongs to --> Cluster

    Deploy|EnvVar
        Deploy --> has many --> EnvVars
        EnvVar --> belongs to --> Deploy