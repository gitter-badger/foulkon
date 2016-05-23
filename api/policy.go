package api

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/database"
)

// Policy domain
type Policy struct {
	ID         string       `json:"ID, omitempty"`
	Name       string       `json:"Name, omitempty"`
	Path       string       `json:"Path, omitempty"`
	Org        string       `json:"Org, omitempty"`
	CreateAt   time.Time    `json:"CreateAt, omitempty"`
	Urn        string       `json:"Urn, omitempty"`
	Statements *[]Statement `json:"Statements, omitempty"`
}

type Statement struct {
	Effect    string   `json:"Effect, omitempty"`
	Action    []string `json:"Action, omitempty"`
	Resources []string `json:"Resources, omitempty"`
}

func (api *AuthAPI) GetPolicy(org string, policyName string) (*Policy, error) {
	// Call repo to retrieve the policy
	policy, err := api.PolicyRepo.GetPolicyByName(org, policyName)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Policy doesn't exist in DB
		if dbError.Code == database.POLICY_NOT_FOUND {
			return nil, &Error{
				Code:    POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Return policy
	return policy, nil
}

func (api *AuthAPI) GetPolicies(org string, pathPrefix string) ([]Policy, error) {
	// Call repo to retrieve the policies
	policies, err := api.PolicyRepo.GetPoliciesFiltered(org, pathPrefix)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Return groups
	return policies, nil
}

func (api *AuthAPI) AddPolicy(name string, path string, org string, statements *[]Statement) (*Policy, error) {
	// Validate fields
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid policy name"),
		}
	}
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid path"),
		}

	}
	if !IsValidStatement(statements) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid statement definition"),
		}

	}

	// Check if policy already exist
	_, err := api.PolicyRepo.GetPolicyByName(org, name)

	// Check if policy could be retrieved
	if err != nil {
		// Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		// Policy doesn't exist in DB
		case database.POLICY_NOT_FOUND:
			// Create policy
			policyCreated, err := api.PolicyRepo.AddPolicy(createPolicy(name, path, org, statements))

			// Check if there is an unexpected error in DB
			if err != nil {
				//Transform to DB error
				dbError := err.(*database.Error)
				return nil, &Error{
					Code:    UNKNOWN_API_ERROR,
					Message: dbError.Message,
				}
			}

			// Return policy created
			return policyCreated, nil
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	} else { // If policy exist it can't create it
		return nil, &Error{
			Code:    POLICY_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create policy, policy with org %v and name %v already exist", org, name),
		}
	}
}

func (api *AuthAPI) UpdatePolicy(org string, policyName string, newName string, newPath string, newStatements []Statement) (*Policy, error) {
	// Call repo to retrieve the policy
	policyDB, err := api.PolicyRepo.GetPolicyByName(org, policyName)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Group doesn't exist in DB
		if dbError.Code == database.POLICY_NOT_FOUND {
			return nil, &Error{
				Code:    POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Validate fields
	if !IsValidName(policyName) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid policy name"),
		}
	}
	if !IsValidPath(newPath) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid path"),
		}

	}
	if !IsValidStatement(&newStatements) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid statement definition"),
		}

	}
	// Check if policy with newName exist
	_, err = api.PolicyRepo.GetPolicyByName(org, newName)

	if err == nil {
		// Policy already exists
		return nil, &Error{
			Code:    POLICY_ALREADY_EXIST,
			Message: fmt.Sprintf("Policy name: %v already exists", newName),
		}
	}

	// Get Urn
	urn := CreateUrn(org, RESOURCE_POLICY, newPath, newName)

	// Update policy
	policy, err := api.PolicyRepo.UpdatePolicy(*policyDB, newName, newPath, urn, newStatements)

	// Check if there is an unexpected error in DB
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	return policy, nil
}

func (api *AuthAPI) DeletePolicy(org string, name string) error {
	// Call repo to retrieve the policy
	policy, err := api.PolicyRepo.GetPolicyByName(org, name)
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		if dbError.Code == database.POLICY_NOT_FOUND {
			return &Error{
				Code:    POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else {
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}

	}

	err = api.PolicyRepo.DeletePolicy(policy.ID)
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Return no error
	return nil
}

func (api *AuthAPI) GetPolicyAttachedGroups(org string, policyName string) ([]Group, error) {
	// Call repo to retrieve the policy
	policy, err := api.PolicyRepo.GetPolicyByName(org, policyName)
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		if dbError.Code == database.POLICY_NOT_FOUND {
			return nil, &Error{
				Code:    POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else {
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}

	}

	// Call repo to retrieve the attached groups
	groups, err := api.PolicyRepo.GetAllPolicyGroupRelation(policy.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Return groups
	return groups, nil
}

func createPolicy(name string, path string, org string, statements *[]Statement) Policy {
	urn := CreateUrn(org, RESOURCE_POLICY, path, name)
	policy := Policy{
		ID:         uuid.NewV4().String(),
		Name:       name,
		Path:       path,
		Org:        org,
		Urn:        urn,
		Statements: statements,
	}

	return policy
}
