package response

type InvalidParameterResponse struct {
	ParamName     string `json:"parameterName"`
	ErrorResponse string `json:"errorMessage"`
}