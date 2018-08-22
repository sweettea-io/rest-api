package buildablesvc

type BuildableLog struct {
  Msg      string
  Level    string
  Complete bool
}

type LogStreamer struct {
  Key     string
  Channel <-chan BuildableLog
}

func NewLogStreamer(buildableUid string) *LogStreamer {
  return &LogStreamer{
    Key: buildableUid,
    Channel: make(chan BuildableLog),
  }
}

func (ls *LogStreamer) GetChannel() <-chan BuildableLog {
  return ls.Channel
}

func (ls *LogStreamer) Watch() {

}
