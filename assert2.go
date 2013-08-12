package bham

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

type aTest struct {
	T       *testing.T
	F       *fTest
	section string
}
type fTest struct {
	T       *testing.T
	section string
}

func (test *aTest) Section(s string) {
	test.section = s
	test.F.section = s
}
func (test *aTest) logDetails() {
	_, fn, l, _ := runtime.Caller(2)
	if test.section != "" {
		test.T.Logf("Error in %s Section in file %s on Line %v",
			test.section,
			fn,
			l,
		)
	} else {
		test.T.Logf("Error in File %s on Line %v", fn, l)
	}
}
func (test *fTest) logDetails() {
	_, fn, l, _ := runtime.Caller(2)
	if test.section != "" {
		test.T.Logf("Error in %s Section in File %s on Line %v",
			test.section,
			fn,
			l,
		)
	} else {
		test.T.Logf("Error in File %s on Line %v", fn, l)
	}
}

// Nil tests
func (test *aTest) IsNil(v interface{}, msgs ...interface{}) {
	if !testIsNil(v) {
		test.logDetails()
		test.T.Log("Nil check failed")
		if len(msgs) > 0 {
			test.T.Error(msgs...)
		} else {
			test.T.Error(v, "is not nil")
		}
	}
}
func (test *aTest) IsNotNil(v interface{}, msgs ...interface{}) {
	if testIsNil(v) {
		test.logDetails()
		test.T.Error(msgs...)
	}
}

// bool tests
func (test *aTest) IsTrue(b bool, msgs ...interface{}) {
	if !b {
		test.logDetails()
		test.T.Error(msgs...)
	}
}
func (test *aTest) IsFalse(b bool, msgs ...interface{}) {
	if b {
		test.logDetails()
		test.T.Error(msgs...)
	}
}

// Equality test
func (test *aTest) AreEqual(x, y interface{}, msgs ...interface{}) {
	if !(reflect.DeepEqual(x, y) || strEqual(x, y)) {
		test.logDetails()
		test.T.Log("Equality check failed")
		if len(msgs) > 0 {
			test.T.Error(msgs...)
		} else {
			test.T.Error(x, "!=", y)
		}
	}
}
func (test *aTest) AreNotEqual(x, y interface{}, msgs ...interface{}) {
	if reflect.DeepEqual(x, y) || strEqual(x, y) {
		test.logDetails()
		test.T.Log("Inequality check failed")
		test.T.Error(msgs...)
	}
}

// String tests
func (test *aTest) StartsWith(s, pre string, msgs ...interface{}) {
	if !strings.HasPrefix(s, pre) {
		test.logDetails()
		test.T.Error(msgs...)
	}
}
func (test *aTest) EndsWith(s, post string, msgs ...interface{}) {
	if !strings.HasSuffix(s, post) {
		test.logDetails()
		test.T.Error(msgs...)
	}
}
func (test *aTest) Matches(s, regex string, msgs ...interface{}) {
	matches, err := regexp.MatchString(regex, s)
	if err != nil {
		test.logDetails()
		panic(err)
	} else if !matches {
		test.logDetails()
		test.T.Error(msgs...)
	}
}
func (test *aTest) NotMatches(s, regex string, msgs ...interface{}) {
	matches, err := regexp.MatchString(regex, s)
	if err != nil {
		test.logDetails()
		panic(err)
	} else if matches {
		test.logDetails()
		test.T.Error(msgs...)
	}
}

// Nil Format tests
func (test *fTest) IsNil(v interface{}, msgFormat string, msgs ...interface{}) {
	if !testIsNil(v) {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *fTest) IsNotNil(v interface{}, msgFormat string, msgs ...interface{}) {
	if testIsNil(v) {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}

// bool tests
func (test *fTest) IsTrue(b bool, msgFormat string, msgs ...interface{}) {
	if !b {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *fTest) IsFalse(b bool, msgFormat string, msgs ...interface{}) {
	if b {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}

// Equality test
func (test *fTest) AreEqual(x, y interface{}, msgFormat string, msgs ...interface{}) {
	if !(reflect.DeepEqual(x, y) || strEqual(x, y)) {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *fTest) AreNotEqual(x, y interface{}, msgFormat string, msgs ...interface{}) {
	if reflect.DeepEqual(x, y) || strEqual(x, y) {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}

// String tests
func (test *fTest) StartsWith(s, pre, msgFormat string, msgs ...interface{}) {
	if !strings.HasPrefix(s, pre) {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *fTest) EndsWith(s, post, msgFormat string, msgs ...interface{}) {
	if !strings.HasSuffix(s, post) {
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *fTest) Matches(s, regex, msgFormat string, msgs ...interface{}) {
	matches, err := regexp.MatchString(regex, s)
	if err != nil {
		test.logDetails()
		panic(err)
	} else if !matches {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *fTest) NotMatches(s, regex, msgFormat string, msgs ...interface{}) {
	matches, err := regexp.MatchString(regex, s)
	if err != nil {
		test.logDetails()
		panic(err)
	} else if matches {
		test.logDetails()
		test.T.Errorf(msgFormat, msgs...)
	}
}

func testIsNil(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}

func strEqual(x, y interface{}) bool {
	return fmt.Sprint(x) == fmt.Sprint(y)
}

func within(t *testing.T, f func(*aTest)) {
	f(&aTest{T: t, F: &fTest{T: t}})
}
