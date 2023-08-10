package matchers

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/in-toto/in-toto-golang/in_toto"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func BeStatementOfType(expectedStatementType, expectedPredicateType string) types.GomegaMatcher {
	return &beStatementOfTypeMatcher{expectedStatementType: expectedStatementType, expectedPredicateType: expectedPredicateType}
}

type beStatementOfTypeMatcher struct {
	expectedStatementType, expectedPredicateType string

	actualStatement *in_toto.Statement
}

func (m *beStatementOfTypeMatcher) Match(actual interface{}) (bool, error) {
	if ok, err := gomega.Not(gomega.BeNil()).Match(actual); !ok || err != nil {
		return ok, err
	}
	switch actual := actual.(type) {
	case in_toto.Statement:
		m.actualStatement = &actual
	case *in_toto.Statement:
		m.actualStatement = actual
	default:
		return false, fmt.Errorf("unexpected object type %T, should be either in_toto.Statement or *in_toto.Statement", actual)
	}
	if ok, err := gomega.Equal(m.expectedStatementType).Match(m.actualStatement.Type); !ok || err != nil {
		return ok, err
	}
	if ok, err := gomega.Equal(m.expectedPredicateType).Match(m.actualStatement.PredicateType); !ok || err != nil {
		return ok, err
	}
	return true, nil
}

func (m *beStatementOfTypeMatcher) FailureMessage(interface{}) string {
	return format.Message(m.actualStatement, fmt.Sprintf("to be a statement of type %q with predicate type %q",
		m.expectedStatementType, m.expectedPredicateType))
}

func (m *beStatementOfTypeMatcher) NegatedFailureMessage(interface{}) string {
	return format.Message(m.actualStatement, fmt.Sprintf("to NOT be a statement of type %q with predicate type %q",
		m.expectedStatementType, m.expectedPredicateType))
}

func HavePredicate(expectedPredicate interface{}) types.GomegaMatcher {
	return &havePredicate{expectedPredicate: expectedPredicate}
}

type havePredicate struct {
	expectedPredicate interface{}

	actualStatement *in_toto.Statement
}

func (m *havePredicate) Match(actual interface{}) (bool, error) {
	if ok, err := gomega.Not(gomega.BeNil()).Match(actual); !ok || err != nil {
		return ok, err
	}
	switch actual := actual.(type) {
	case in_toto.Statement:
		m.actualStatement = &actual
	case *in_toto.Statement:
		m.actualStatement = actual
	default:
		return false, fmt.Errorf("unexpected object type %T, should be either in_toto.Statement or *in_toto.Statement", actual)
	}
	actualPredicateData, err := json.Marshal(m.actualStatement.Predicate)
	if err != nil {
		return false, err
	}

	actualPredicate := new(interface{})
	if err := json.Unmarshal(actualPredicateData, actualPredicate); err != nil {
		return false, err
	}

	if ok, err := gomega.Equal(m.expectedPredicate).Match(actualPredicate); !ok || err != nil {
		return ok, err
	}
	return true, nil
}

func (m *havePredicate) FailureMessage(interface{}) string {
	return format.Message(m.actualStatement.Predicate, "to be", m.expectedPredicate)
}

func (m *havePredicate) NegatedFailureMessage(interface{}) string {
	return format.Message(m.actualStatement.Predicate, "to NOT be", m.expectedPredicate)

}

func HavePredicateOfTypeSatisfying(expectedPredicateType interface{}, callback func(interface{})) *havePredicateOfTypeSatisfying {
	return &havePredicateOfTypeSatisfying{expectedPredicateType: expectedPredicateType, callback: callback}
}

type havePredicateOfTypeSatisfying struct {
	expectedPredicateType interface{}
	callback              func(interface{})

	actualStatement *in_toto.Statement
}

func (m *havePredicateOfTypeSatisfying) Match(actual interface{}) (bool, error) {
	if ok, err := gomega.Not(gomega.BeNil()).Match(actual); !ok || err != nil {
		return ok, err
	}
	switch actual := actual.(type) {
	case in_toto.Statement:
		m.actualStatement = &actual
	case *in_toto.Statement:
		m.actualStatement = actual
	default:
		return false, fmt.Errorf("unexpected object type %T, should be either in_toto.Statement or *in_toto.Statement", actual)
	}
	actualPredicate := reflect.New(reflect.TypeOf(m.expectedPredicateType))

	actualPredicateData, err := json.Marshal(m.actualStatement.Predicate)
	if err != nil {
		return false, err
	}

	actualPredicateObj := actualPredicate.Interface()

	if err := json.Unmarshal(actualPredicateData, actualPredicateObj); err != nil {
		return false, fmt.Errorf("cannot unmarshal %s: %w", string(actualPredicateData), err)
	}

	m.callback(actualPredicateObj)
	return true, nil
}

func (m *havePredicateOfTypeSatisfying) FailureMessage(interface{}) string {
	return format.Message(m.actualStatement.Predicate, "to satisfy custom test via callback")
}

func (m *havePredicateOfTypeSatisfying) NegatedFailureMessage(interface{}) string {
	return format.Message(m.actualStatement.Predicate, "to NOT satisfy custom test via callback")
}
