package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
)

func TestWorkerHandler_HandleGetUsers(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		pathPrefix string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetUserExternalIDsResponse
		expectedError      api.Error
		// Manager Results
		getListUsersResult []string
		// Manager Errors
		getListUsersErr error
	}{
		"OkCase": {
			pathPrefix:         "myPath",
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetUserExternalIDsResponse{
				ExternalIDs: []string{"userId1", "userId2"},
			},
			getListUsersResult: []string{"userId1", "userId2"},
		},
		"ErrorCaseUnauthorizedError": {
			pathPrefix:         "myPath",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Error",
			},
			getListUsersErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Error",
			},
		},
		"ErrorCaseUnknownApiError": {
			expectedStatusCode: http.StatusInternalServerError,
			getListUsersErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetListUsersMethod][0] = test.getListUsersResult
		testApi.ArgsOut[GetListUsersMethod][1] = test.getListUsersErr

		req, err := http.NewRequest(http.MethodGet, server.URL+USER_ROOT_URL, nil)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		if test.pathPrefix != "" {
			q := req.URL.Query()
			q.Add("PathPrefix", test.pathPrefix)
			req.URL.RawQuery = q.Encode()
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		// Check received parameter
		if testApi.ArgsIn[GetListUsersMethod][1] != test.pathPrefix {
			t.Errorf("Test case %v. Received different PathPrefix (wanted:%v / received:%v)", n, test.pathPrefix, testApi.ArgsIn[GetListUsersMethod][1])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			getUserExternalIDsResponse := GetUserExternalIDsResponse{}
			err = json.NewDecoder(res.Body).Decode(&getUserExternalIDsResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(getUserExternalIDsResponse, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v",
					n, diff)
				continue
			}
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing error response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(apiError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v",
					n, diff)
				continue
			}

		}

	}
}

func TestWorkerHandler_HandlePostUsers(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		request *CreateUserRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   CreateUserResponse
		expectedError      api.Error
		// Manager Results
		addUserResult *api.User
		// Manager Errors
		addUserErr error
	}{
		"OkCase": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: CreateUserResponse{
				User: &api.User{
					ID:         "UserID",
					ExternalID: "ExternalID",
					Path:       "Path",
					Urn:        "urn",
					CreateAt:   now,
				},
			},
			addUserResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseUserAlreadyExist": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.USER_ALREADY_EXIST,
				Message: "User already exist",
			},
			addUserErr: &api.Error{
				Code:    api.USER_ALREADY_EXIST,
				Message: "User already exist",
			},
		},
		"ErrorCaseInvalidParameterError": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			addUserErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			addUserErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusInternalServerError,
			addUserErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[AddUserMethod][0] = test.addUserResult
		testApi.ArgsOut[AddUserMethod][1] = test.addUserErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			if err != nil {
				t.Errorf("Test case %v. Unexpected marshalling api request %v", n, err)
				continue
			}
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		req, err := http.NewRequest(http.MethodPost, server.URL+USER_ROOT_URL, body)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		if test.request != nil {
			// Check received parameters
			if testApi.ArgsIn[AddUserMethod][1] != test.request.ExternalID {
				t.Errorf("Test case %v. Received different ExternalID (wanted:%v / received:%v)", n, test.request.ExternalID, testApi.ArgsIn[AddUserMethod][1])
				continue
			}
			if testApi.ArgsIn[AddUserMethod][2] != test.request.Path {
				t.Errorf("Test case %v. Received different Path (wanted:%v / received:%v)", n, test.request.Path, testApi.ArgsIn[AddUserMethod][2])
				continue
			}
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusCreated:
			createUserResponse := CreateUserResponse{}
			err = json.NewDecoder(res.Body).Decode(&createUserResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(createUserResponse, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v",
					n, diff)
				continue
			}
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing error response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(apiError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v",
					n, diff)
				continue
			}

		}

	}
}
