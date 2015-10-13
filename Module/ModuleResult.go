package Module

type ModuleResult struct {
	Event       interface{}
	ReturnCode  int
	LogMessages struct {
					Timestamp string
					Level     string
					Message   string
					Source    string
				}
}