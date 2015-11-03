package ConfigLayouts

//EventsPerTime file layout
type EventsPerTime struct {
	InfluxDB struct {
		Server   string
		Database string
		Username string
		Password string
	}
}
