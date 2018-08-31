package mocks

import "github.com/sweettea-io/rest-api/internal/app"

type MockConfig struct {
  MockBaseRoute func() string
  MockOnTest func() bool
  MockOnLocal func() bool
  MockOnDev func() bool
  MockOnStaging func() bool
  MockOnProd func() bool
  MockTrainClusterConfigured func() bool
  MockBuildpackEnvs func() map[string]string
}

func (mc *MockConfig) BaseRoute() string {
  if mc.MockBaseRoute != nil {
    return mc.MockBaseRoute()
  }

  return app.Config.BaseRoute()
}

func (mc *MockConfig) OnTest() bool {
  if mc.MockOnTest != nil {
    return mc.MockOnTest()
  }

  return app.Config.OnTest()
}

func (mc *MockConfig) OnLocal() bool {
  if mc.MockOnLocal != nil {
    return mc.MockOnLocal()
  }

  return app.Config.OnLocal()
}

func (mc *MockConfig) OnDev() bool {
  if mc.MockOnDev != nil {
    return mc.MockOnDev()
  }

  return app.Config.OnDev()
}

func (mc *MockConfig) OnStaging() bool {
  if mc.MockOnStaging != nil {
    return mc.MockOnStaging()
  }

  return app.Config.OnStaging()
}

func (mc *MockConfig) OnProd() bool {
  if mc.MockOnProd != nil {
    return mc.MockOnProd()
  }

  return app.Config.OnProd()
}

func (mc *MockConfig) TrainClusterConfigured() bool {
  if mc.MockTrainClusterConfigured != nil {
    return mc.MockTrainClusterConfigured()
  }

  return app.Config.TrainClusterConfigured()
}

func (mc *MockConfig) BuildpackEnvs() map[string]string {
  if mc.MockBuildpackEnvs != nil {
    return mc.MockBuildpackEnvs()
  }

  return app.Config.BuildpackEnvs()
}