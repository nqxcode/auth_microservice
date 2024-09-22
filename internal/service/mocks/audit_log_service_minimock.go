// Code generated by http://github.com/gojuno/minimock (v3.3.14). DO NOT EDIT.

package mocks

//go:generate minimock -i github.com/nqxcode/auth_microservice/internal/service.AuditLogService -o audit_log_service_minimock.go -n AuditLogServiceMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/nqxcode/auth_microservice/internal/model"
)

// AuditLogServiceMock implements service.AuditLogService
type AuditLogServiceMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcCreate          func(ctx context.Context, message *model.Log) (err error)
	inspectFuncCreate   func(ctx context.Context, message *model.Log)
	afterCreateCounter  uint64
	beforeCreateCounter uint64
	CreateMock          mAuditLogServiceMockCreate
}

// NewAuditLogServiceMock returns a mock for service.AuditLogService
func NewAuditLogServiceMock(t minimock.Tester) *AuditLogServiceMock {
	m := &AuditLogServiceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateMock = mAuditLogServiceMockCreate{mock: m}
	m.CreateMock.callArgs = []*AuditLogServiceMockCreateParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mAuditLogServiceMockCreate struct {
	optional           bool
	mock               *AuditLogServiceMock
	defaultExpectation *AuditLogServiceMockCreateExpectation
	expectations       []*AuditLogServiceMockCreateExpectation

	callArgs []*AuditLogServiceMockCreateParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// AuditLogServiceMockCreateExpectation specifies expectation struct of the AuditLogService.Create
type AuditLogServiceMockCreateExpectation struct {
	mock      *AuditLogServiceMock
	params    *AuditLogServiceMockCreateParams
	paramPtrs *AuditLogServiceMockCreateParamPtrs
	results   *AuditLogServiceMockCreateResults
	Counter   uint64
}

// AuditLogServiceMockCreateParams contains parameters of the AuditLogService.Create
type AuditLogServiceMockCreateParams struct {
	ctx     context.Context
	message *model.Log
}

// AuditLogServiceMockCreateParamPtrs contains pointers to parameters of the AuditLogService.Create
type AuditLogServiceMockCreateParamPtrs struct {
	ctx     *context.Context
	message **model.Log
}

// AuditLogServiceMockCreateResults contains results of the AuditLogService.Create
type AuditLogServiceMockCreateResults struct {
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmCreate *mAuditLogServiceMockCreate) Optional() *mAuditLogServiceMockCreate {
	mmCreate.optional = true
	return mmCreate
}

// Expect sets up expected params for AuditLogService.Create
func (mmCreate *mAuditLogServiceMockCreate) Expect(ctx context.Context, message *model.Log) *mAuditLogServiceMockCreate {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("AuditLogServiceMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &AuditLogServiceMockCreateExpectation{}
	}

	if mmCreate.defaultExpectation.paramPtrs != nil {
		mmCreate.mock.t.Fatalf("AuditLogServiceMock.Create mock is already set by ExpectParams functions")
	}

	mmCreate.defaultExpectation.params = &AuditLogServiceMockCreateParams{ctx, message}
	for _, e := range mmCreate.expectations {
		if minimock.Equal(e.params, mmCreate.defaultExpectation.params) {
			mmCreate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmCreate.defaultExpectation.params)
		}
	}

	return mmCreate
}

// ExpectCtxParam1 sets up expected param ctx for AuditLogService.Create
func (mmCreate *mAuditLogServiceMockCreate) ExpectCtxParam1(ctx context.Context) *mAuditLogServiceMockCreate {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("AuditLogServiceMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &AuditLogServiceMockCreateExpectation{}
	}

	if mmCreate.defaultExpectation.params != nil {
		mmCreate.mock.t.Fatalf("AuditLogServiceMock.Create mock is already set by Expect")
	}

	if mmCreate.defaultExpectation.paramPtrs == nil {
		mmCreate.defaultExpectation.paramPtrs = &AuditLogServiceMockCreateParamPtrs{}
	}
	mmCreate.defaultExpectation.paramPtrs.ctx = &ctx

	return mmCreate
}

// ExpectMessageParam2 sets up expected param message for AuditLogService.Create
func (mmCreate *mAuditLogServiceMockCreate) ExpectMessageParam2(message *model.Log) *mAuditLogServiceMockCreate {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("AuditLogServiceMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &AuditLogServiceMockCreateExpectation{}
	}

	if mmCreate.defaultExpectation.params != nil {
		mmCreate.mock.t.Fatalf("AuditLogServiceMock.Create mock is already set by Expect")
	}

	if mmCreate.defaultExpectation.paramPtrs == nil {
		mmCreate.defaultExpectation.paramPtrs = &AuditLogServiceMockCreateParamPtrs{}
	}
	mmCreate.defaultExpectation.paramPtrs.message = &message

	return mmCreate
}

// Inspect accepts an inspector function that has same arguments as the AuditLogService.Create
func (mmCreate *mAuditLogServiceMockCreate) Inspect(f func(ctx context.Context, message *model.Log)) *mAuditLogServiceMockCreate {
	if mmCreate.mock.inspectFuncCreate != nil {
		mmCreate.mock.t.Fatalf("Inspect function is already set for AuditLogServiceMock.Create")
	}

	mmCreate.mock.inspectFuncCreate = f

	return mmCreate
}

// Return sets up results that will be returned by AuditLogService.Create
func (mmCreate *mAuditLogServiceMockCreate) Return(err error) *AuditLogServiceMock {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("AuditLogServiceMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &AuditLogServiceMockCreateExpectation{mock: mmCreate.mock}
	}
	mmCreate.defaultExpectation.results = &AuditLogServiceMockCreateResults{err}
	return mmCreate.mock
}

// Set uses given function f to mock the AuditLogService.Create method
func (mmCreate *mAuditLogServiceMockCreate) Set(f func(ctx context.Context, message *model.Log) (err error)) *AuditLogServiceMock {
	if mmCreate.defaultExpectation != nil {
		mmCreate.mock.t.Fatalf("Default expectation is already set for the AuditLogService.Create method")
	}

	if len(mmCreate.expectations) > 0 {
		mmCreate.mock.t.Fatalf("Some expectations are already set for the AuditLogService.Create method")
	}

	mmCreate.mock.funcCreate = f
	return mmCreate.mock
}

// When sets expectation for the AuditLogService.Create which will trigger the result defined by the following
// Then helper
func (mmCreate *mAuditLogServiceMockCreate) When(ctx context.Context, message *model.Log) *AuditLogServiceMockCreateExpectation {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("AuditLogServiceMock.Create mock is already set by Set")
	}

	expectation := &AuditLogServiceMockCreateExpectation{
		mock:   mmCreate.mock,
		params: &AuditLogServiceMockCreateParams{ctx, message},
	}
	mmCreate.expectations = append(mmCreate.expectations, expectation)
	return expectation
}

// Then sets up AuditLogService.Create return parameters for the expectation previously defined by the When method
func (e *AuditLogServiceMockCreateExpectation) Then(err error) *AuditLogServiceMock {
	e.results = &AuditLogServiceMockCreateResults{err}
	return e.mock
}

// Times sets number of times AuditLogService.Create should be invoked
func (mmCreate *mAuditLogServiceMockCreate) Times(n uint64) *mAuditLogServiceMockCreate {
	if n == 0 {
		mmCreate.mock.t.Fatalf("Times of AuditLogServiceMock.Create mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmCreate.expectedInvocations, n)
	return mmCreate
}

func (mmCreate *mAuditLogServiceMockCreate) invocationsDone() bool {
	if len(mmCreate.expectations) == 0 && mmCreate.defaultExpectation == nil && mmCreate.mock.funcCreate == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmCreate.mock.afterCreateCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmCreate.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Create implements service.AuditLogService
func (mmCreate *AuditLogServiceMock) Create(ctx context.Context, message *model.Log) (err error) {
	mm_atomic.AddUint64(&mmCreate.beforeCreateCounter, 1)
	defer mm_atomic.AddUint64(&mmCreate.afterCreateCounter, 1)

	if mmCreate.inspectFuncCreate != nil {
		mmCreate.inspectFuncCreate(ctx, message)
	}

	mm_params := AuditLogServiceMockCreateParams{ctx, message}

	// Record call args
	mmCreate.CreateMock.mutex.Lock()
	mmCreate.CreateMock.callArgs = append(mmCreate.CreateMock.callArgs, &mm_params)
	mmCreate.CreateMock.mutex.Unlock()

	for _, e := range mmCreate.CreateMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmCreate.CreateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmCreate.CreateMock.defaultExpectation.Counter, 1)
		mm_want := mmCreate.CreateMock.defaultExpectation.params
		mm_want_ptrs := mmCreate.CreateMock.defaultExpectation.paramPtrs

		mm_got := AuditLogServiceMockCreateParams{ctx, message}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmCreate.t.Errorf("AuditLogServiceMock.Create got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.message != nil && !minimock.Equal(*mm_want_ptrs.message, mm_got.message) {
				mmCreate.t.Errorf("AuditLogServiceMock.Create got unexpected parameter message, want: %#v, got: %#v%s\n", *mm_want_ptrs.message, mm_got.message, minimock.Diff(*mm_want_ptrs.message, mm_got.message))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmCreate.t.Errorf("AuditLogServiceMock.Create got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmCreate.CreateMock.defaultExpectation.results
		if mm_results == nil {
			mmCreate.t.Fatal("No results are set for the AuditLogServiceMock.Create")
		}
		return (*mm_results).err
	}
	if mmCreate.funcCreate != nil {
		return mmCreate.funcCreate(ctx, message)
	}
	mmCreate.t.Fatalf("Unexpected call to AuditLogServiceMock.Create. %v %v", ctx, message)
	return
}

// CreateAfterCounter returns a count of finished AuditLogServiceMock.Create invocations
func (mmCreate *AuditLogServiceMock) CreateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCreate.afterCreateCounter)
}

// CreateBeforeCounter returns a count of AuditLogServiceMock.Create invocations
func (mmCreate *AuditLogServiceMock) CreateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCreate.beforeCreateCounter)
}

// Calls returns a list of arguments used in each call to AuditLogServiceMock.Create.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmCreate *mAuditLogServiceMockCreate) Calls() []*AuditLogServiceMockCreateParams {
	mmCreate.mutex.RLock()

	argCopy := make([]*AuditLogServiceMockCreateParams, len(mmCreate.callArgs))
	copy(argCopy, mmCreate.callArgs)

	mmCreate.mutex.RUnlock()

	return argCopy
}

// MinimockCreateDone returns true if the count of the Create invocations corresponds
// the number of defined expectations
func (m *AuditLogServiceMock) MinimockCreateDone() bool {
	if m.CreateMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.CreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.CreateMock.invocationsDone()
}

// MinimockCreateInspect logs each unmet expectation
func (m *AuditLogServiceMock) MinimockCreateInspect() {
	for _, e := range m.CreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AuditLogServiceMock.Create with params: %#v", *e.params)
		}
	}

	afterCreateCounter := mm_atomic.LoadUint64(&m.afterCreateCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.CreateMock.defaultExpectation != nil && afterCreateCounter < 1 {
		if m.CreateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to AuditLogServiceMock.Create")
		} else {
			m.t.Errorf("Expected call to AuditLogServiceMock.Create with params: %#v", *m.CreateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCreate != nil && afterCreateCounter < 1 {
		m.t.Error("Expected call to AuditLogServiceMock.Create")
	}

	if !m.CreateMock.invocationsDone() && afterCreateCounter > 0 {
		m.t.Errorf("Expected %d calls to AuditLogServiceMock.Create but found %d calls",
			mm_atomic.LoadUint64(&m.CreateMock.expectedInvocations), afterCreateCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *AuditLogServiceMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockCreateInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *AuditLogServiceMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *AuditLogServiceMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockCreateDone()
}
