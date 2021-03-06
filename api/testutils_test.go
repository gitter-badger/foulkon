package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/kylelemons/godebug/pretty"
	"math/rand"
	"testing"
)

const (
	GetUserByExternalIDMethod = "GetUserByExternalID"
	AddUserMethod             = "AddUser"
	UpdateUserMethod          = "UpdateUser"
	GetUsersFilteredMethod    = "GetUsersFiltered"
	GetGroupsByUserIDMethod   = "GetGroupsByUserID"
	RemoveUserMethod          = "RemoveUser"
	GetGroupByNameMethod      = "GetGroupByName"
	IsMemberOfGroupMethod     = "IsMemberOfGroup"
	GetGroupMembersMethod     = "GetGroupMembers"
	IsAttachedToGroupMethod   = "IsAttachedToGroup"
	GetAttachedPoliciesMethod = "GetAttachedPolicies"
	GetGroupsFilteredMethod   = "GetGroupsFiltered"
	RemoveGroupMethod         = "RemoveGroup"
	AddGroupMethod            = "AddGroup"
	AddMemberMethod           = "AddMember"
	RemoveMemberMethod        = "RemoveMember"
	UpdateGroupMethod         = "UpdateGroup"
	AttachPolicyMethod        = "AttachPolicy"
	DetachPolicyMethod        = "DetachPolicy"
	GetPolicyByNameMethod     = "GetPolicyByName"
	AddPolicyMethod           = "AddPolicy"
	UpdatePolicyMethod        = "UpdatePolicy"
	RemovePolicyMethod        = "RemovePolicy"
	GetPoliciesFilteredMethod = "GetPoliciesFiltered"
	GetAttachedGroupsMethod   = "GetAttachedGroups"
)

// TestRepo that implements all repo manager interfaces
type TestRepo struct {
	ArgsIn       map[string][]interface{}
	ArgsOut      map[string][]interface{}
	SpecialFuncs map[string]interface{}
}

// func that initializes the TestRepo
func makeTestRepo() *TestRepo {
	testRepo := &TestRepo{
		ArgsIn:       make(map[string][]interface{}),
		ArgsOut:      make(map[string][]interface{}),
		SpecialFuncs: make(map[string]interface{}),
	}
	testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[AddUserMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[UpdateUserMethod] = make([]interface{}, 3)
	testRepo.ArgsIn[GetUsersFilteredMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[GetGroupsByUserIDMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[RemoveUserMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[GetGroupByNameMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[IsMemberOfGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[GetGroupMembersMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[IsAttachedToGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[GetAttachedPoliciesMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[GetGroupsFilteredMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[RemoveGroupMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[AddGroupMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[AddMemberMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[RemoveMemberMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[UpdateGroupMethod] = make([]interface{}, 4)
	testRepo.ArgsIn[AttachPolicyMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[DetachPolicyMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[GetPolicyByNameMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[AddPolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[UpdatePolicyMethod] = make([]interface{}, 5)
	testRepo.ArgsIn[RemovePolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[GetPoliciesFilteredMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[GetAttachedGroupsMethod] = make([]interface{}, 1)

	testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[AddUserMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[UpdateUserMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetUsersFilteredMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetGroupsByUserIDMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[RemoveUserMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[GetGroupByNameMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[IsMemberOfGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetGroupMembersMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[IsAttachedToGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetAttachedPoliciesMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetGroupsFilteredMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[RemoveGroupMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[AddGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[AddMemberMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[RemoveMemberMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[UpdateGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[AttachPolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[DetachPolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[GetPolicyByNameMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[AddPolicyMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[UpdatePolicyMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[RemovePolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[GetPoliciesFilteredMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetAttachedGroupsMethod] = make([]interface{}, 2)

	return testRepo
}

func makeTestAPI(testRepo *TestRepo) *AuthAPI {
	api := &AuthAPI{
		UserRepo:   testRepo,
		GroupRepo:  testRepo,
		PolicyRepo: testRepo,
		Logger:     logrus.StandardLogger(),
	}
	return api
}

//////////////////
// User repo
//////////////////

func (t TestRepo) GetUserByExternalID(id string) (*User, error) {
	t.ArgsIn[GetUserByExternalIDMethod][0] = id
	if specialFunc, ok := t.SpecialFuncs[GetUserByExternalIDMethod].(func(id string) (*User, error)); ok && specialFunc != nil {
		return specialFunc(id)
	}
	var user *User
	if t.ArgsOut[GetUserByExternalIDMethod][0] != nil {
		user = t.ArgsOut[GetUserByExternalIDMethod][0].(*User)
	}
	var err error
	if t.ArgsOut[GetUserByExternalIDMethod][1] != nil {
		err = t.ArgsOut[GetUserByExternalIDMethod][1].(error)
	}
	return user, err
}

func (t TestRepo) AddUser(user User) (*User, error) {
	t.ArgsIn[AddUserMethod][0] = user
	var created *User
	if t.ArgsOut[AddUserMethod][0] != nil {
		created = t.ArgsOut[AddUserMethod][0].(*User)
	}
	var err error
	if t.ArgsOut[AddUserMethod][1] != nil {
		err = t.ArgsOut[AddUserMethod][1].(error)
	}
	return created, err
}

func (t TestRepo) UpdateUser(user User, newPath string, newUrn string) (*User, error) {
	t.ArgsIn[UpdateUserMethod][0] = user
	t.ArgsIn[UpdateUserMethod][1] = newPath
	t.ArgsIn[UpdateUserMethod][2] = newUrn
	var updated *User
	if t.ArgsOut[UpdateUserMethod][0] != nil {
		updated = t.ArgsOut[UpdateUserMethod][0].(*User)
	}
	var err error
	if t.ArgsOut[UpdateUserMethod][1] != nil {
		err = t.ArgsOut[UpdateUserMethod][1].(error)
	}
	return updated, err
}

func (t TestRepo) GetUsersFiltered(pathPrefix string) ([]User, error) {
	t.ArgsIn[GetUsersFilteredMethod][0] = pathPrefix
	var users []User
	if t.ArgsOut[GetUsersFilteredMethod][0] != nil {
		users = t.ArgsOut[GetUsersFilteredMethod][0].([]User)
	}
	var err error
	if t.ArgsOut[GetUsersFilteredMethod][1] != nil {
		err = t.ArgsOut[GetUsersFilteredMethod][1].(error)
	}
	return users, err
}

func (t TestRepo) GetGroupsByUserID(id string) ([]Group, error) {
	t.ArgsIn[GetGroupsByUserIDMethod][0] = id
	var groups []Group
	if t.ArgsOut[GetGroupsByUserIDMethod][0] != nil {
		groups = t.ArgsOut[GetGroupsByUserIDMethod][0].([]Group)
	}
	var err error
	if t.ArgsOut[GetGroupsByUserIDMethod][1] != nil {
		err = t.ArgsOut[GetGroupsByUserIDMethod][1].(error)
	}
	return groups, err
}

func (t TestRepo) RemoveUser(id string) error {
	t.ArgsIn[RemoveUserMethod][0] = id
	var err error
	if t.ArgsOut[RemoveUserMethod][0] != nil {
		err = t.ArgsOut[RemoveUserMethod][0].(error)
	}
	return err
}

//////////////////
// Group repo
//////////////////

func (t TestRepo) GetGroupByName(org string, name string) (*Group, error) {
	t.ArgsIn[GetGroupByNameMethod][0] = org
	t.ArgsIn[GetGroupByNameMethod][1] = name
	if specialFunc, ok := t.SpecialFuncs[GetGroupByNameMethod].(func(org string, name string) (*Group, error)); ok && specialFunc != nil {
		return specialFunc(org, name)
	}
	var group *Group
	if t.ArgsOut[GetGroupByNameMethod][0] != nil {
		group = t.ArgsOut[GetGroupByNameMethod][0].(*Group)
	}
	var err error
	if t.ArgsOut[GetGroupByNameMethod][1] != nil {
		err = t.ArgsOut[GetGroupByNameMethod][1].(error)
	}
	return group, err
}

func (t TestRepo) IsMemberOfGroup(userID string, groupID string) (bool, error) {
	t.ArgsIn[IsMemberOfGroupMethod][0] = userID
	t.ArgsIn[IsMemberOfGroupMethod][1] = groupID
	var isMember bool
	if t.ArgsOut[IsMemberOfGroupMethod][0] != nil {
		isMember = t.ArgsOut[IsMemberOfGroupMethod][0].(bool)
	}
	var err error
	if t.ArgsOut[IsMemberOfGroupMethod][1] != nil {
		err = t.ArgsOut[IsMemberOfGroupMethod][1].(error)
	}
	return isMember, err
}

func (t TestRepo) GetGroupMembers(groupID string) ([]User, error) {
	t.ArgsIn[GetGroupMembersMethod][0] = groupID
	var members []User
	if t.ArgsOut[GetGroupMembersMethod][0] != nil {
		members = t.ArgsOut[GetGroupMembersMethod][0].([]User)
	}
	var err error
	if t.ArgsOut[GetGroupMembersMethod][1] != nil {
		err = t.ArgsOut[GetGroupMembersMethod][1].(error)
	}
	return members, err
}

func (t TestRepo) IsAttachedToGroup(groupID string, policyID string) (bool, error) {
	t.ArgsIn[IsAttachedToGroupMethod][0] = groupID
	t.ArgsIn[IsAttachedToGroupMethod][1] = policyID
	var isAttached bool
	if t.ArgsOut[IsAttachedToGroupMethod][0] != nil {
		isAttached = t.ArgsOut[IsAttachedToGroupMethod][0].(bool)
	}
	var err error
	if t.ArgsOut[IsAttachedToGroupMethod][1] != nil {
		err = t.ArgsOut[IsAttachedToGroupMethod][1].(error)
	}
	return isAttached, err
}

func (t TestRepo) GetAttachedPolicies(groupID string) ([]Policy, error) {
	t.ArgsIn[GetAttachedPoliciesMethod][0] = groupID
	var policies []Policy
	if t.ArgsOut[GetAttachedPoliciesMethod][0] != nil {
		policies = t.ArgsOut[GetAttachedPoliciesMethod][0].([]Policy)
	}
	var err error
	if t.ArgsOut[GetAttachedPoliciesMethod][1] != nil {
		err = t.ArgsOut[GetAttachedPoliciesMethod][1].(error)
	}
	return policies, err
}

func (t TestRepo) GetGroupsFiltered(org string, pathPrefix string) ([]Group, error) {
	t.ArgsIn[GetGroupsFilteredMethod][0] = org
	t.ArgsIn[GetGroupsFilteredMethod][1] = pathPrefix

	var groups []Group
	if t.ArgsOut[GetGroupsFilteredMethod][0] != nil {
		groups = t.ArgsOut[GetGroupsFilteredMethod][0].([]Group)
	}
	var err error
	if t.ArgsOut[GetGroupsFilteredMethod][1] != nil {
		err = t.ArgsOut[GetGroupsFilteredMethod][1].(error)
	}
	return groups, err
}
func (t TestRepo) RemoveGroup(id string) error {
	t.ArgsIn[RemoveGroupMethod][0] = id
	var err error
	if t.ArgsOut[RemoveGroupMethod][0] != nil {
		err = t.ArgsOut[RemoveGroupMethod][0].(error)
	}
	return err
}

func (t TestRepo) AddGroup(group Group) (*Group, error) {
	t.ArgsIn[AddGroupMethod][0] = group
	var created *Group
	if t.ArgsOut[AddGroupMethod][0] != nil {
		created = t.ArgsOut[AddGroupMethod][0].(*Group)
	}
	var err error
	if t.ArgsOut[AddGroupMethod][1] != nil {
		err = t.ArgsOut[AddGroupMethod][1].(error)
	}
	return created, err
}

func (t TestRepo) AddMember(userID string, groupID string) error {
	t.ArgsIn[AddMemberMethod][0] = userID
	t.ArgsIn[AddMemberMethod][1] = groupID
	var err error
	if t.ArgsOut[AddMemberMethod][0] != nil {
		err = t.ArgsOut[AddMemberMethod][0].(error)
	}
	return err
}

func (t TestRepo) RemoveMember(userID string, groupID string) error {
	t.ArgsIn[RemoveMemberMethod][0] = userID
	t.ArgsIn[RemoveMemberMethod][1] = groupID
	var err error
	if t.ArgsOut[RemoveMemberMethod][0] != nil {
		err = t.ArgsOut[RemoveMemberMethod][0].(error)
	}
	return err
}

func (t TestRepo) UpdateGroup(group Group, newName string, newPath string, newUrn string) (*Group, error) {
	t.ArgsIn[UpdateGroupMethod][0] = group
	t.ArgsIn[UpdateGroupMethod][1] = newName
	t.ArgsIn[UpdateGroupMethod][2] = newPath
	t.ArgsIn[UpdateGroupMethod][3] = newUrn

	var updated *Group
	if t.ArgsOut[UpdateGroupMethod][0] != nil {
		updated = t.ArgsOut[UpdateGroupMethod][0].(*Group)
	}
	var err error
	if t.ArgsOut[UpdateGroupMethod][1] != nil {
		err = t.ArgsOut[UpdateGroupMethod][1].(error)
	}
	return updated, err
}

func (t TestRepo) AttachPolicy(groupID string, policyID string) error {
	t.ArgsIn[AttachPolicyMethod][0] = groupID
	t.ArgsIn[AttachPolicyMethod][1] = policyID
	var err error
	if t.ArgsOut[AttachPolicyMethod][0] != nil {
		err = t.ArgsOut[AttachPolicyMethod][0].(error)
	}
	return err
}
func (t TestRepo) DetachPolicy(groupID string, policyID string) error {
	t.ArgsIn[DetachPolicyMethod][0] = groupID
	t.ArgsIn[DetachPolicyMethod][1] = policyID
	var err error
	if t.ArgsOut[DetachPolicyMethod][0] != nil {
		err = t.ArgsOut[DetachPolicyMethod][0].(error)
	}
	return err
}

//////////////////
// Policy repo
//////////////////

func (t TestRepo) GetPolicyByName(org string, name string) (*Policy, error) {
	t.ArgsIn[GetPolicyByNameMethod][0] = org
	t.ArgsIn[GetPolicyByNameMethod][1] = name
	if specialFunc, ok := t.SpecialFuncs[GetPolicyByNameMethod].(func(org string, name string) (*Policy, error)); ok && specialFunc != nil {
		return specialFunc(org, name)
	}
	var policy *Policy
	if t.ArgsOut[GetPolicyByNameMethod][0] != nil {
		policy = t.ArgsOut[GetPolicyByNameMethod][0].(*Policy)
	}
	var err error
	if t.ArgsOut[GetPolicyByNameMethod][1] != nil {
		err = t.ArgsOut[GetPolicyByNameMethod][1].(error)
	}
	return policy, err
}

func (t TestRepo) AddPolicy(policy Policy) (*Policy, error) {
	t.ArgsIn[AddPolicyMethod][0] = policy
	var created *Policy
	if t.ArgsOut[AddPolicyMethod][0] != nil {
		created = t.ArgsOut[AddPolicyMethod][0].(*Policy)
	}
	var err error
	if t.ArgsOut[AddPolicyMethod][1] != nil {
		err = t.ArgsOut[AddPolicyMethod][1].(error)
	}
	return created, err
}

func (t TestRepo) UpdatePolicy(policy Policy, newName string, newPath string, newUrn string, newStatements []Statement) (*Policy, error) {
	t.ArgsIn[UpdatePolicyMethod][0] = policy
	t.ArgsIn[UpdatePolicyMethod][1] = newName
	t.ArgsIn[UpdatePolicyMethod][2] = newPath
	t.ArgsIn[UpdatePolicyMethod][3] = newUrn
	t.ArgsIn[UpdatePolicyMethod][4] = newStatements

	var updated *Policy
	if t.ArgsOut[UpdatePolicyMethod][0] != nil {
		updated = t.ArgsOut[UpdatePolicyMethod][0].(*Policy)
	}
	var err error
	if t.ArgsOut[UpdatePolicyMethod][1] != nil {
		err = t.ArgsOut[UpdatePolicyMethod][1].(error)
	}
	return updated, err
}

func (t TestRepo) RemovePolicy(id string) error {
	t.ArgsIn[RemovePolicyMethod][0] = id
	var err error
	if t.ArgsOut[RemovePolicyMethod][0] != nil {
		err = t.ArgsOut[RemovePolicyMethod][0].(error)
	}
	return err
}

func (t TestRepo) GetPoliciesFiltered(org string, pathPrefix string) ([]Policy, error) {
	t.ArgsIn[GetPoliciesFilteredMethod][0] = org
	t.ArgsIn[GetPoliciesFilteredMethod][1] = pathPrefix

	var policies []Policy
	if t.ArgsOut[GetPoliciesFilteredMethod][0] != nil {
		policies = t.ArgsOut[GetPoliciesFilteredMethod][0].([]Policy)
	}
	var err error
	if t.ArgsOut[GetPoliciesFilteredMethod][1] != nil {
		err = t.ArgsOut[GetPoliciesFilteredMethod][1].(error)
	}
	return policies, err
}

func (t TestRepo) GetAttachedGroups(policyID string) ([]Group, error) {
	t.ArgsIn[GetAttachedGroupsMethod][0] = policyID

	var groups []Group
	if t.ArgsOut[GetAttachedGroupsMethod][0] != nil {
		groups = t.ArgsOut[GetAttachedGroupsMethod][0].([]Group)
	}
	var err error
	if t.ArgsOut[GetAttachedGroupsMethod][1] != nil {
		err = t.ArgsOut[GetAttachedGroupsMethod][1].(error)
	}
	return groups, err
}

// Private helper methods

func GetRandomString(runeValue []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runeValue[rand.Intn(len(runeValue))]
	}
	return string(b)
}

func checkMethodResponse(t *testing.T, testcase string, expectedError error, receivedError error, expectedResponse interface{}, receivedResponse interface{}) {
	if expectedError != nil {
		apiError, ok := receivedError.(*Error)
		if !ok || apiError == nil {
			t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", testcase, receivedError)
			return
		}
		if diff := pretty.Compare(apiError, expectedError); diff != "" {
			t.Errorf("Test %v failed. Received different errors (received/wanted) %v", testcase, diff)
			return
		}
	} else {
		if receivedError != nil {
			t.Errorf("Test %v failed: %v", testcase, receivedError)
			return
		} else {
			if diff := pretty.Compare(receivedResponse, expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", testcase, diff)
				return
			}
		}
	}
}
