package auth

import "github.com/casbin/govaluate"

type IEnforcer interface {
	GetAllSubjects() ([]string, error)
	GetAllNamedSubjects(ptype string) ([]string, error)
	GetAllObjects() ([]string, error)
	GetAllNamedObjects(ptype string) ([]string, error)
	GetAllActions() ([]string, error)
	GetAllNamedActions(ptype string) ([]string, error)
	GetAllRoles() ([]string, error)
	GetAllNamedRoles(ptype string) ([]string, error)
	GetPolicy() ([][]string, error)
	GetFilteredPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)
	GetNamedPolicy(ptype string) ([][]string, error)
	GetFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)
	GetGroupingPolicy() ([][]string, error)
	GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)
	GetNamedGroupingPolicy(ptype string) ([][]string, error)
	GetFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)
	GetFilteredNamedPolicyWithMatcher(ptype string, matcher string) ([][]string, error)
	HasPolicy(params ...interface{}) (bool, error)
	HasNamedPolicy(ptype string, params ...interface{}) (bool, error)
	AddPolicy(params ...interface{}) (bool, error)
	AddPolicies(rules [][]string) (bool, error)
	AddPoliciesEx(rules [][]string) (bool, error)
	AddNamedPolicy(ptype string, params ...interface{}) (bool, error)
	AddNamedPolicies(ptype string, rules [][]string) (bool, error)
	AddNamedPoliciesEx(ptype string, rules [][]string) (bool, error)
	RemovePolicy(params ...interface{}) (bool, error)
	UpdatePolicy(oldPolicy []string, newPolicy []string) (bool, error)
	UpdateNamedPolicy(ptype string, p1 []string, p2 []string) (bool, error)
	UpdatePolicies(oldPolices [][]string, newPolicies [][]string) (bool, error)
	UpdateNamedPolicies(ptype string, p1 [][]string, p2 [][]string) (bool, error)
	UpdateFilteredPolicies(newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)
	UpdateFilteredNamedPolicies(ptype string, newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)
	RemovePolicies(rules [][]string) (bool, error)
	RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error)
	RemoveNamedPolicy(ptype string, params ...interface{}) (bool, error)
	RemoveNamedPolicies(ptype string, rules [][]string) (bool, error)
	RemoveFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)
	HasGroupingPolicy(params ...interface{}) (bool, error)
	HasNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
	AddGroupingPolicy(params ...interface{}) (bool, error)
	AddGroupingPolicies(rules [][]string) (bool, error)
	AddGroupingPoliciesEx(rules [][]string) (bool, error)
	AddNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
	AddNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)
	AddNamedGroupingPoliciesEx(ptype string, rules [][]string) (bool, error)
	RemoveGroupingPolicy(params ...interface{}) (bool, error)
	RemoveGroupingPolicies(rules [][]string) (bool, error)
	RemoveFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) (bool, error)
	RemoveNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
	RemoveNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)
	UpdateGroupingPolicy(oldRule []string, newRule []string) (bool, error)
	UpdateGroupingPolicies(oldRules [][]string, newRules [][]string) (bool, error)
	UpdateNamedGroupingPolicy(ptype string, oldRule []string, newRule []string) (bool, error)
	UpdateNamedGroupingPolicies(ptype string, oldRules [][]string, newRules [][]string) (bool, error)
	RemoveFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)
	AddFunction(name string, function govaluate.ExpressionFunction)
	SelfAddPolicy(sec string, ptype string, rule []string) (bool, error)
	SelfAddPolicies(sec string, ptype string, rules [][]string) (bool, error)
	SelfAddPoliciesEx(sec string, ptype string, rules [][]string) (bool, error)
	SelfRemovePolicy(sec string, ptype string, rule []string) (bool, error)
	SelfRemovePolicies(sec string, ptype string, rules [][]string) (bool, error)
	SelfRemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error)
	SelfUpdatePolicy(sec string, ptype string, oldRule, newRule []string) (bool, error)
	SelfUpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) (bool, error)
}

type EnforcerEmpty struct {
}

func NewEnforcerEmpty() IEnforcer {
	return &EnforcerEmpty{}
}
func (e *EnforcerEmpty) GetAllSubjects() ([]string, error) {
	return []string{}, nil
}

func (e *EnforcerEmpty) GetAllNamedSubjects(ptype string) ([]string, error) {
	return []string{}, nil
}

func (e *EnforcerEmpty) GetAllObjects() ([]string, error) {
	return []string{}, nil
}

func (e *EnforcerEmpty) GetAllNamedObjects(ptype string) ([]string, error) {
	return []string{}, nil
}

func (e *EnforcerEmpty) GetAllActions() ([]string, error) {
	return []string{}, nil
}

func (e *EnforcerEmpty) GetAllNamedActions(ptype string) ([]string, error) {
	return []string{}, nil
}

func (e *EnforcerEmpty) GetAllRoles() ([]string, error) {
	return []string{}, nil
}

func (e *EnforcerEmpty) GetAllNamedRoles(ptype string) ([]string, error) {
	return []string{}, nil
}

func (e *EnforcerEmpty) GetPolicy() ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) GetFilteredPolicy(fieldIndex int, fieldValues ...string) ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) GetNamedPolicy(ptype string) ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) GetFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) GetGroupingPolicy() ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) GetNamedGroupingPolicy(ptype string) ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) GetFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) GetFilteredNamedPolicyWithMatcher(ptype string, matcher string) ([][]string, error) {
	return [][]string{}, nil
}

func (e *EnforcerEmpty) HasPolicy(params ...interface{}) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) HasNamedPolicy(ptype string, params ...interface{}) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) AddPolicy(params ...interface{}) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) AddPolicies(rules [][]string) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) AddPoliciesEx(rules [][]string) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) AddNamedPolicy(ptype string, params ...interface{}) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) AddNamedPolicies(ptype string, rules [][]string) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) AddNamedPoliciesEx(ptype string, rules [][]string) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) RemovePolicy(params ...interface{}) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdatePolicy(oldPolicy []string, newPolicy []string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdateNamedPolicy(ptype string, p1 []string, p2 []string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdatePolicies(oldPolices [][]string, newPolicies [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdateNamedPolicies(ptype string, p1 [][]string, p2 [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdateFilteredPolicies(newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdateFilteredNamedPolicies(ptype string, newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemovePolicies(rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveNamedPolicy(ptype string, params ...interface{}) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveNamedPolicies(ptype string, rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) HasGroupingPolicy(params ...interface{}) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) HasNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error) {
	return false, nil
}

func (e *EnforcerEmpty) AddGroupingPolicy(params ...interface{}) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) AddGroupingPolicies(rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) AddGroupingPoliciesEx(rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) AddNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) AddNamedGroupingPolicies(ptype string, rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) AddNamedGroupingPoliciesEx(ptype string, rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveGroupingPolicy(params ...interface{}) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveGroupingPolicies(rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveNamedGroupingPolicies(ptype string, rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdateGroupingPolicy(oldRule []string, newRule []string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdateGroupingPolicies(oldRules [][]string, newRules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdateNamedGroupingPolicy(ptype string, oldRule []string, newRule []string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) UpdateNamedGroupingPolicies(ptype string, oldRules [][]string, newRules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) RemoveFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) AddFunction(name string, function govaluate.ExpressionFunction) {

}

func (e *EnforcerEmpty) SelfAddPolicy(sec string, ptype string, rule []string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) SelfAddPolicies(sec string, ptype string, rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) SelfAddPoliciesEx(sec string, ptype string, rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) SelfRemovePolicy(sec string, ptype string, rule []string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) SelfRemovePolicies(sec string, ptype string, rules [][]string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) SelfRemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) SelfUpdatePolicy(sec string, ptype string, oldRule, newRule []string) (bool, error) {
	return true, nil
}

func (e *EnforcerEmpty) SelfUpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) (bool, error) {
	return true, nil
}
