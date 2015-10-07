package Config

type Config struct {
	RuleSystem struct {
				   Enabled    bool
				   ModulePath string
				   Rulefile   string
				   Worker     int
			   }
}