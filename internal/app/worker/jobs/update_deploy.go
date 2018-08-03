package jobs

// UpdateDeploy needs to call the right next job based on what actually needs to be updated. Options include:
// (1) sha --> new image needs to be built and then PATCHED on the deployment
// (2) model/modelVersion --> env vars need to be PATCHED on the deployment