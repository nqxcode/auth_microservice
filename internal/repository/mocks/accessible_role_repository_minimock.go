// Code generated by http://github.com/gojuno/minimock (v3.3.14). DO NOT EDIT.

package mocks

//go:generate minimock -i github.com/nqxcode/auth_microservice/internal/repository.AccessibleRoleRepository -o accessible_role_repository_minimock.go -n AccessibleRoleRepositoryMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/nqxcode/auth_microservice/internal/model"
)

// AccessibleRoleRepositoryMock implements repository.AccessibleRoleRepository
type AccessibleRoleRepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetList          func(ctx context.Context) (aa1 []model.AccessibleRole, err error)
	inspectFuncGetList   func(ctx context.Context)
	afterGetListCounter  uint64
	beforeGetListCounter uint64
	GetListMock          mAccessibleRoleRepositoryMockGetList
}

// NewAccessibleRoleRepositoryMock returns a mock for repository.AccessibleRoleRepository
func NewAccessibleRoleRepositoryMock(t minimock.Tester) *AccessibleRoleRepositoryMock {
	m := &AccessibleRoleRepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetListMock = mAccessibleRoleRepositoryMockGetList{mock: m}
	m.GetListMock.callArgs = []*AccessibleRoleRepositoryMockGetListParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mAccessibleRoleRepositoryMockGetList struct {
	optional           bool
	mock               *AccessibleRoleRepositoryMock
	defaultExpectation *AccessibleRoleRepositoryMockGetListExpectation
	expectations       []*AccessibleRoleRepositoryMockGetListExpectation

	callArgs []*AccessibleRoleRepositoryMockGetListParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// AccessibleRoleRepositoryMockGetListExpectation specifies expectation struct of the AccessibleRoleRepository.GetList
type AccessibleRoleRepositoryMockGetListExpectation struct {
	mock      *AccessibleRoleRepositoryMock
	params    *AccessibleRoleRepositoryMockGetListParams
	paramPtrs *AccessibleRoleRepositoryMockGetListParamPtrs
	results   *AccessibleRoleRepositoryMockGetListResults
	Counter   uint64
}

// AccessibleRoleRepositoryMockGetListParams contains parameters of the AccessibleRoleRepository.GetList
type AccessibleRoleRepositoryMockGetListParams struct {
	ctx context.Context
}

// AccessibleRoleRepositoryMockGetListParamPtrs contains pointers to parameters of the AccessibleRoleRepository.GetList
type AccessibleRoleRepositoryMockGetListParamPtrs struct {
	ctx *context.Context
}

// AccessibleRoleRepositoryMockGetListResults contains results of the AccessibleRoleRepository.GetList
type AccessibleRoleRepositoryMockGetListResults struct {
	aa1 []model.AccessibleRole
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetList *mAccessibleRoleRepositoryMockGetList) Optional() *mAccessibleRoleRepositoryMockGetList {
	mmGetList.optional = true
	return mmGetList
}

// Expect sets up expected params for AccessibleRoleRepository.GetList
func (mmGetList *mAccessibleRoleRepositoryMockGetList) Expect(ctx context.Context) *mAccessibleRoleRepositoryMockGetList {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AccessibleRoleRepositoryMock.GetList mock is already set by Set")
	}

	if mmGetList.defaultExpectation == nil {
		mmGetList.defaultExpectation = &AccessibleRoleRepositoryMockGetListExpectation{}
	}

	if mmGetList.defaultExpectation.paramPtrs != nil {
		mmGetList.mock.t.Fatalf("AccessibleRoleRepositoryMock.GetList mock is already set by ExpectParams functions")
	}

	mmGetList.defaultExpectation.params = &AccessibleRoleRepositoryMockGetListParams{ctx}
	for _, e := range mmGetList.expectations {
		if minimock.Equal(e.params, mmGetList.defaultExpectation.params) {
			mmGetList.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetList.defaultExpectation.params)
		}
	}

	return mmGetList
}

// ExpectCtxParam1 sets up expected param ctx for AccessibleRoleRepository.GetList
func (mmGetList *mAccessibleRoleRepositoryMockGetList) ExpectCtxParam1(ctx context.Context) *mAccessibleRoleRepositoryMockGetList {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AccessibleRoleRepositoryMock.GetList mock is already set by Set")
	}

	if mmGetList.defaultExpectation == nil {
		mmGetList.defaultExpectation = &AccessibleRoleRepositoryMockGetListExpectation{}
	}

	if mmGetList.defaultExpectation.params != nil {
		mmGetList.mock.t.Fatalf("AccessibleRoleRepositoryMock.GetList mock is already set by Expect")
	}

	if mmGetList.defaultExpectation.paramPtrs == nil {
		mmGetList.defaultExpectation.paramPtrs = &AccessibleRoleRepositoryMockGetListParamPtrs{}
	}
	mmGetList.defaultExpectation.paramPtrs.ctx = &ctx

	return mmGetList
}

// Inspect accepts an inspector function that has same arguments as the AccessibleRoleRepository.GetList
func (mmGetList *mAccessibleRoleRepositoryMockGetList) Inspect(f func(ctx context.Context)) *mAccessibleRoleRepositoryMockGetList {
	if mmGetList.mock.inspectFuncGetList != nil {
		mmGetList.mock.t.Fatalf("Inspect function is already set for AccessibleRoleRepositoryMock.GetList")
	}

	mmGetList.mock.inspectFuncGetList = f

	return mmGetList
}

// Return sets up results that will be returned by AccessibleRoleRepository.GetList
func (mmGetList *mAccessibleRoleRepositoryMockGetList) Return(aa1 []model.AccessibleRole, err error) *AccessibleRoleRepositoryMock {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AccessibleRoleRepositoryMock.GetList mock is already set by Set")
	}

	if mmGetList.defaultExpectation == nil {
		mmGetList.defaultExpectation = &AccessibleRoleRepositoryMockGetListExpectation{mock: mmGetList.mock}
	}
	mmGetList.defaultExpectation.results = &AccessibleRoleRepositoryMockGetListResults{aa1, err}
	return mmGetList.mock
}

// Set uses given function f to mock the AccessibleRoleRepository.GetList method
func (mmGetList *mAccessibleRoleRepositoryMockGetList) Set(f func(ctx context.Context) (aa1 []model.AccessibleRole, err error)) *AccessibleRoleRepositoryMock {
	if mmGetList.defaultExpectation != nil {
		mmGetList.mock.t.Fatalf("Default expectation is already set for the AccessibleRoleRepository.GetList method")
	}

	if len(mmGetList.expectations) > 0 {
		mmGetList.mock.t.Fatalf("Some expectations are already set for the AccessibleRoleRepository.GetList method")
	}

	mmGetList.mock.funcGetList = f
	return mmGetList.mock
}

// When sets expectation for the AccessibleRoleRepository.GetList which will trigger the result defined by the following
// Then helper
func (mmGetList *mAccessibleRoleRepositoryMockGetList) When(ctx context.Context) *AccessibleRoleRepositoryMockGetListExpectation {
	if mmGetList.mock.funcGetList != nil {
		mmGetList.mock.t.Fatalf("AccessibleRoleRepositoryMock.GetList mock is already set by Set")
	}

	expectation := &AccessibleRoleRepositoryMockGetListExpectation{
		mock:   mmGetList.mock,
		params: &AccessibleRoleRepositoryMockGetListParams{ctx},
	}
	mmGetList.expectations = append(mmGetList.expectations, expectation)
	return expectation
}

// Then sets up AccessibleRoleRepository.GetList return parameters for the expectation previously defined by the When method
func (e *AccessibleRoleRepositoryMockGetListExpectation) Then(aa1 []model.AccessibleRole, err error) *AccessibleRoleRepositoryMock {
	e.results = &AccessibleRoleRepositoryMockGetListResults{aa1, err}
	return e.mock
}

// Times sets number of times AccessibleRoleRepository.GetList should be invoked
func (mmGetList *mAccessibleRoleRepositoryMockGetList) Times(n uint64) *mAccessibleRoleRepositoryMockGetList {
	if n == 0 {
		mmGetList.mock.t.Fatalf("Times of AccessibleRoleRepositoryMock.GetList mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetList.expectedInvocations, n)
	return mmGetList
}

func (mmGetList *mAccessibleRoleRepositoryMockGetList) invocationsDone() bool {
	if len(mmGetList.expectations) == 0 && mmGetList.defaultExpectation == nil && mmGetList.mock.funcGetList == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetList.mock.afterGetListCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetList.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetList implements repository.AccessibleRoleRepository
func (mmGetList *AccessibleRoleRepositoryMock) GetList(ctx context.Context) (aa1 []model.AccessibleRole, err error) {
	mm_atomic.AddUint64(&mmGetList.beforeGetListCounter, 1)
	defer mm_atomic.AddUint64(&mmGetList.afterGetListCounter, 1)

	if mmGetList.inspectFuncGetList != nil {
		mmGetList.inspectFuncGetList(ctx)
	}

	mm_params := AccessibleRoleRepositoryMockGetListParams{ctx}

	// Record call args
	mmGetList.GetListMock.mutex.Lock()
	mmGetList.GetListMock.callArgs = append(mmGetList.GetListMock.callArgs, &mm_params)
	mmGetList.GetListMock.mutex.Unlock()

	for _, e := range mmGetList.GetListMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.aa1, e.results.err
		}
	}

	if mmGetList.GetListMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetList.GetListMock.defaultExpectation.Counter, 1)
		mm_want := mmGetList.GetListMock.defaultExpectation.params
		mm_want_ptrs := mmGetList.GetListMock.defaultExpectation.paramPtrs

		mm_got := AccessibleRoleRepositoryMockGetListParams{ctx}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetList.t.Errorf("AccessibleRoleRepositoryMock.GetList got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetList.t.Errorf("AccessibleRoleRepositoryMock.GetList got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetList.GetListMock.defaultExpectation.results
		if mm_results == nil {
			mmGetList.t.Fatal("No results are set for the AccessibleRoleRepositoryMock.GetList")
		}
		return (*mm_results).aa1, (*mm_results).err
	}
	if mmGetList.funcGetList != nil {
		return mmGetList.funcGetList(ctx)
	}
	mmGetList.t.Fatalf("Unexpected call to AccessibleRoleRepositoryMock.GetList. %v", ctx)
	return
}

// GetListAfterCounter returns a count of finished AccessibleRoleRepositoryMock.GetList invocations
func (mmGetList *AccessibleRoleRepositoryMock) GetListAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetList.afterGetListCounter)
}

// GetListBeforeCounter returns a count of AccessibleRoleRepositoryMock.GetList invocations
func (mmGetList *AccessibleRoleRepositoryMock) GetListBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetList.beforeGetListCounter)
}

// Calls returns a list of arguments used in each call to AccessibleRoleRepositoryMock.GetList.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetList *mAccessibleRoleRepositoryMockGetList) Calls() []*AccessibleRoleRepositoryMockGetListParams {
	mmGetList.mutex.RLock()

	argCopy := make([]*AccessibleRoleRepositoryMockGetListParams, len(mmGetList.callArgs))
	copy(argCopy, mmGetList.callArgs)

	mmGetList.mutex.RUnlock()

	return argCopy
}

// MinimockGetListDone returns true if the count of the GetList invocations corresponds
// the number of defined expectations
func (m *AccessibleRoleRepositoryMock) MinimockGetListDone() bool {
	if m.GetListMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetListMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetListMock.invocationsDone()
}

// MinimockGetListInspect logs each unmet expectation
func (m *AccessibleRoleRepositoryMock) MinimockGetListInspect() {
	for _, e := range m.GetListMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AccessibleRoleRepositoryMock.GetList with params: %#v", *e.params)
		}
	}

	afterGetListCounter := mm_atomic.LoadUint64(&m.afterGetListCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetListMock.defaultExpectation != nil && afterGetListCounter < 1 {
		if m.GetListMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to AccessibleRoleRepositoryMock.GetList")
		} else {
			m.t.Errorf("Expected call to AccessibleRoleRepositoryMock.GetList with params: %#v", *m.GetListMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetList != nil && afterGetListCounter < 1 {
		m.t.Error("Expected call to AccessibleRoleRepositoryMock.GetList")
	}

	if !m.GetListMock.invocationsDone() && afterGetListCounter > 0 {
		m.t.Errorf("Expected %d calls to AccessibleRoleRepositoryMock.GetList but found %d calls",
			mm_atomic.LoadUint64(&m.GetListMock.expectedInvocations), afterGetListCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *AccessibleRoleRepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetListInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *AccessibleRoleRepositoryMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *AccessibleRoleRepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetListDone()
}
