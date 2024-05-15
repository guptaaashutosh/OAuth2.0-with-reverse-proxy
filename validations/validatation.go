package validation

import (
	"fmt"
	"learn/httpserver/response"
	"net/http"

	// "github.com/gookit/validate"
	"github.com/guptaaashutosh/go_validate"
)

// ValidateParameters validates queryParameters, requestParameters and urlParameters
func ValidateParameters(r *http.Request, requestBodyDto interface{}, requestParametersMap *map[string]string, requestParametersFiltersMap *map[string]string, queryParameterMap *map[string]string, queryParametersFiltersMap *map[string]string, urlParamsError *[]response.InvalidParameterResponse) (error, error, error) {
	var invalidParamErr []response.InvalidParameterResponse
	var invalidParamKeysArr []string
	

	if urlParamsError != nil {
		invalidParamErr = append(invalidParamErr, *urlParamsError...)
	}
	// Collection of query parameter
	queryParameterData := validate.FromURLValues(r.URL.Query())
	queryParamData := queryParameterData.Create()

	var queryParameterMapDereference map[string]string
	if queryParameterMap != nil {
		queryParameterMapDereference = *queryParameterMap
		queryParamData.StringRules(queryParameterMapDereference)
	}

	if queryParametersFiltersMap != nil {
		queryParamData.FilterRules(*queryParametersFiltersMap)
	}

	if !queryParamData.Validate() {
		// Range over query parameter map to get particular parameter error
		for key, _ := range queryParameterMapDereference {
			if len(queryParamData.Errors.FieldOne(key)) != 0 {
				invalidParamErr = append(invalidParamErr, response.InvalidParameterResponse{ParamName: key, ErrorResponse: queryParamData.Errors.FieldOne(key)})
				invalidParamKeysArr = append(invalidParamKeysArr, key)
			}
		}
	}

	// Collection of request parameter
	requestParameterData, err := validate.FromRequest(r)
	if err != nil {
		// invalidRequestError := er.GenerateError(http.StatusBadRequest, err.Error())
		return err, nil, nil
	}

	requestBodyData := requestParameterData.Create()

	var requestParametersMapDereference map[string]string
	if requestParametersMap != nil {
		requestParametersMapDereference = *requestParametersMap
		requestBodyData.StringRules(requestParametersMapDereference)
	}

	if requestParametersFiltersMap != nil {
		requestBodyData.FilterRules(*requestParametersFiltersMap)
	}

	if !requestBodyData.Validate() {
		// Range over request parameter map to get particular parameter error
		for key, _ := range requestParametersMapDereference {
			if len(requestBodyData.Errors.FieldOne(key)) != 0 {
				invalidParamErr = append(invalidParamErr, response.InvalidParameterResponse{ParamName: key, ErrorResponse: requestBodyData.Errors.FieldOne(key)})
				invalidParamKeysArr = append(invalidParamKeysArr, key)
			}
		}
	}
	var httpErrorCode int
	if len(invalidParamErr) == 0 {
		if httpErrorCode, err = queryParamData.BindSafeData(requestBodyDto); err != nil {
			return err, nil, nil
		}
		//printing to resolve unused error 
		fmt.Print(httpErrorCode)
		
		if httpErrorCode, err = requestBodyData.BindSafeData(requestBodyDto); err != nil {
			return err, nil, nil
		}
		//printing to resolve unused error 
		fmt.Print(httpErrorCode)
	}

	// invalidParamsSingleLineErrMsg := er.GenerateError(http.StatusBadRequest, fmt.Sprintf("Invalid data in:%s", strings.Join(invalidParamKeysArr, ", ")))
	return nil, err, err
}
