package singleutils

import "testing"

type ITestService interface {
	GetName() string
}

type TestService struct {
	name string
}

func (s *TestService) GetName() string {
	return s.name
}

func Test_Get(t *testing.T) {
	/*	if err := Set[ITestService](func() ITestService {
			v := &TestService{}
			v.name = "ITestService"
			return v
		}); err != nil {
			t.Error(err)
		}
		if s, err := Get[ITestService](); err != nil {
			t.Error(err)
		} else {
			t.Logf("s.name = %s", s.GetName())
		}*/

	if err := Set[TestService](func() TestService { return TestService{} }); err != nil {
		t.Error(err)
	}
	if s, err := Get[TestService](); err != nil {
		t.Error(err)
	} else {
		s.name = "TestService"
		t.Logf("s.name = %s", s.name)
	}

	if err := Set[*TestService](func() *TestService { return &TestService{} }); err != nil {
		t.Error(err)
	}
	if s, err := Get[*TestService](); err != nil {
		t.Error(err)
	} else {
		s.name = "*TestService"
		t.Logf("s.name = %s", s.name)
	}

}
