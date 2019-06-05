package resources

type TrackerError struct {
	Code             string `json:"code"`
	Kind             string `json:"kind"`
	Error            string `json:"error"`
	GeneralProblem   string `json:"general_problem"`
	ValidationErrors []struct {
		Field   string `json:"field"`
		Problem string `json:"problem"`
	} `json:"validation_errors"`
}

func (t *TrackerError) Validate() bool {
	return t.Code != ""
}