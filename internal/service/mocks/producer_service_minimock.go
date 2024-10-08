// Code generated by http://github.com/gojuno/minimock (v3.3.14). DO NOT EDIT.

package mocks

//go:generate minimock -i github.com/nqxcode/auth_microservice/internal/service.ProducerService -o producer_service_minimock.go -n ProducerServiceMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/nqxcode/auth_microservice/internal/model"
)

// ProducerServiceMock implements service.ProducerService
type ProducerServiceMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcSendMessage          func(ctx context.Context, message model.LogMessage) (err error)
	inspectFuncSendMessage   func(ctx context.Context, message model.LogMessage)
	afterSendMessageCounter  uint64
	beforeSendMessageCounter uint64
	SendMessageMock          mProducerServiceMockSendMessage
}

// NewProducerServiceMock returns a mock for service.ProducerService
func NewProducerServiceMock(t minimock.Tester) *ProducerServiceMock {
	m := &ProducerServiceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SendMessageMock = mProducerServiceMockSendMessage{mock: m}
	m.SendMessageMock.callArgs = []*ProducerServiceMockSendMessageParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mProducerServiceMockSendMessage struct {
	optional           bool
	mock               *ProducerServiceMock
	defaultExpectation *ProducerServiceMockSendMessageExpectation
	expectations       []*ProducerServiceMockSendMessageExpectation

	callArgs []*ProducerServiceMockSendMessageParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// ProducerServiceMockSendMessageExpectation specifies expectation struct of the ProducerService.SendMessage
type ProducerServiceMockSendMessageExpectation struct {
	mock      *ProducerServiceMock
	params    *ProducerServiceMockSendMessageParams
	paramPtrs *ProducerServiceMockSendMessageParamPtrs
	results   *ProducerServiceMockSendMessageResults
	Counter   uint64
}

// ProducerServiceMockSendMessageParams contains parameters of the ProducerService.SendMessage
type ProducerServiceMockSendMessageParams struct {
	ctx     context.Context
	message model.LogMessage
}

// ProducerServiceMockSendMessageParamPtrs contains pointers to parameters of the ProducerService.SendMessage
type ProducerServiceMockSendMessageParamPtrs struct {
	ctx     *context.Context
	message *model.LogMessage
}

// ProducerServiceMockSendMessageResults contains results of the ProducerService.SendMessage
type ProducerServiceMockSendMessageResults struct {
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmSendMessage *mProducerServiceMockSendMessage) Optional() *mProducerServiceMockSendMessage {
	mmSendMessage.optional = true
	return mmSendMessage
}

// Expect sets up expected params for ProducerService.SendMessage
func (mmSendMessage *mProducerServiceMockSendMessage) Expect(ctx context.Context, message model.LogMessage) *mProducerServiceMockSendMessage {
	if mmSendMessage.mock.funcSendMessage != nil {
		mmSendMessage.mock.t.Fatalf("ProducerServiceMock.SendMessage mock is already set by Set")
	}

	if mmSendMessage.defaultExpectation == nil {
		mmSendMessage.defaultExpectation = &ProducerServiceMockSendMessageExpectation{}
	}

	if mmSendMessage.defaultExpectation.paramPtrs != nil {
		mmSendMessage.mock.t.Fatalf("ProducerServiceMock.SendMessage mock is already set by ExpectParams functions")
	}

	mmSendMessage.defaultExpectation.params = &ProducerServiceMockSendMessageParams{ctx, message}
	for _, e := range mmSendMessage.expectations {
		if minimock.Equal(e.params, mmSendMessage.defaultExpectation.params) {
			mmSendMessage.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSendMessage.defaultExpectation.params)
		}
	}

	return mmSendMessage
}

// ExpectCtxParam1 sets up expected param ctx for ProducerService.SendMessage
func (mmSendMessage *mProducerServiceMockSendMessage) ExpectCtxParam1(ctx context.Context) *mProducerServiceMockSendMessage {
	if mmSendMessage.mock.funcSendMessage != nil {
		mmSendMessage.mock.t.Fatalf("ProducerServiceMock.SendMessage mock is already set by Set")
	}

	if mmSendMessage.defaultExpectation == nil {
		mmSendMessage.defaultExpectation = &ProducerServiceMockSendMessageExpectation{}
	}

	if mmSendMessage.defaultExpectation.params != nil {
		mmSendMessage.mock.t.Fatalf("ProducerServiceMock.SendMessage mock is already set by Expect")
	}

	if mmSendMessage.defaultExpectation.paramPtrs == nil {
		mmSendMessage.defaultExpectation.paramPtrs = &ProducerServiceMockSendMessageParamPtrs{}
	}
	mmSendMessage.defaultExpectation.paramPtrs.ctx = &ctx

	return mmSendMessage
}

// ExpectMessageParam2 sets up expected param message for ProducerService.SendMessage
func (mmSendMessage *mProducerServiceMockSendMessage) ExpectMessageParam2(message model.LogMessage) *mProducerServiceMockSendMessage {
	if mmSendMessage.mock.funcSendMessage != nil {
		mmSendMessage.mock.t.Fatalf("ProducerServiceMock.SendMessage mock is already set by Set")
	}

	if mmSendMessage.defaultExpectation == nil {
		mmSendMessage.defaultExpectation = &ProducerServiceMockSendMessageExpectation{}
	}

	if mmSendMessage.defaultExpectation.params != nil {
		mmSendMessage.mock.t.Fatalf("ProducerServiceMock.SendMessage mock is already set by Expect")
	}

	if mmSendMessage.defaultExpectation.paramPtrs == nil {
		mmSendMessage.defaultExpectation.paramPtrs = &ProducerServiceMockSendMessageParamPtrs{}
	}
	mmSendMessage.defaultExpectation.paramPtrs.message = &message

	return mmSendMessage
}

// Inspect accepts an inspector function that has same arguments as the ProducerService.SendMessage
func (mmSendMessage *mProducerServiceMockSendMessage) Inspect(f func(ctx context.Context, message model.LogMessage)) *mProducerServiceMockSendMessage {
	if mmSendMessage.mock.inspectFuncSendMessage != nil {
		mmSendMessage.mock.t.Fatalf("Inspect function is already set for ProducerServiceMock.SendMessage")
	}

	mmSendMessage.mock.inspectFuncSendMessage = f

	return mmSendMessage
}

// Return sets up results that will be returned by ProducerService.SendMessage
func (mmSendMessage *mProducerServiceMockSendMessage) Return(err error) *ProducerServiceMock {
	if mmSendMessage.mock.funcSendMessage != nil {
		mmSendMessage.mock.t.Fatalf("ProducerServiceMock.SendMessage mock is already set by Set")
	}

	if mmSendMessage.defaultExpectation == nil {
		mmSendMessage.defaultExpectation = &ProducerServiceMockSendMessageExpectation{mock: mmSendMessage.mock}
	}
	mmSendMessage.defaultExpectation.results = &ProducerServiceMockSendMessageResults{err}
	return mmSendMessage.mock
}

// Set uses given function f to mock the ProducerService.SendMessage method
func (mmSendMessage *mProducerServiceMockSendMessage) Set(f func(ctx context.Context, message model.LogMessage) (err error)) *ProducerServiceMock {
	if mmSendMessage.defaultExpectation != nil {
		mmSendMessage.mock.t.Fatalf("Default expectation is already set for the ProducerService.SendMessage method")
	}

	if len(mmSendMessage.expectations) > 0 {
		mmSendMessage.mock.t.Fatalf("Some expectations are already set for the ProducerService.SendMessage method")
	}

	mmSendMessage.mock.funcSendMessage = f
	return mmSendMessage.mock
}

// When sets expectation for the ProducerService.SendMessage which will trigger the result defined by the following
// Then helper
func (mmSendMessage *mProducerServiceMockSendMessage) When(ctx context.Context, message model.LogMessage) *ProducerServiceMockSendMessageExpectation {
	if mmSendMessage.mock.funcSendMessage != nil {
		mmSendMessage.mock.t.Fatalf("ProducerServiceMock.SendMessage mock is already set by Set")
	}

	expectation := &ProducerServiceMockSendMessageExpectation{
		mock:   mmSendMessage.mock,
		params: &ProducerServiceMockSendMessageParams{ctx, message},
	}
	mmSendMessage.expectations = append(mmSendMessage.expectations, expectation)
	return expectation
}

// Then sets up ProducerService.SendMessage return parameters for the expectation previously defined by the When method
func (e *ProducerServiceMockSendMessageExpectation) Then(err error) *ProducerServiceMock {
	e.results = &ProducerServiceMockSendMessageResults{err}
	return e.mock
}

// Times sets number of times ProducerService.SendMessage should be invoked
func (mmSendMessage *mProducerServiceMockSendMessage) Times(n uint64) *mProducerServiceMockSendMessage {
	if n == 0 {
		mmSendMessage.mock.t.Fatalf("Times of ProducerServiceMock.SendMessage mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmSendMessage.expectedInvocations, n)
	return mmSendMessage
}

func (mmSendMessage *mProducerServiceMockSendMessage) invocationsDone() bool {
	if len(mmSendMessage.expectations) == 0 && mmSendMessage.defaultExpectation == nil && mmSendMessage.mock.funcSendMessage == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmSendMessage.mock.afterSendMessageCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmSendMessage.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// SendMessage implements service.ProducerService
func (mmSendMessage *ProducerServiceMock) SendMessage(ctx context.Context, message model.LogMessage) (err error) {
	mm_atomic.AddUint64(&mmSendMessage.beforeSendMessageCounter, 1)
	defer mm_atomic.AddUint64(&mmSendMessage.afterSendMessageCounter, 1)

	if mmSendMessage.inspectFuncSendMessage != nil {
		mmSendMessage.inspectFuncSendMessage(ctx, message)
	}

	mm_params := ProducerServiceMockSendMessageParams{ctx, message}

	// Record call args
	mmSendMessage.SendMessageMock.mutex.Lock()
	mmSendMessage.SendMessageMock.callArgs = append(mmSendMessage.SendMessageMock.callArgs, &mm_params)
	mmSendMessage.SendMessageMock.mutex.Unlock()

	for _, e := range mmSendMessage.SendMessageMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmSendMessage.SendMessageMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSendMessage.SendMessageMock.defaultExpectation.Counter, 1)
		mm_want := mmSendMessage.SendMessageMock.defaultExpectation.params
		mm_want_ptrs := mmSendMessage.SendMessageMock.defaultExpectation.paramPtrs

		mm_got := ProducerServiceMockSendMessageParams{ctx, message}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmSendMessage.t.Errorf("ProducerServiceMock.SendMessage got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.message != nil && !minimock.Equal(*mm_want_ptrs.message, mm_got.message) {
				mmSendMessage.t.Errorf("ProducerServiceMock.SendMessage got unexpected parameter message, want: %#v, got: %#v%s\n", *mm_want_ptrs.message, mm_got.message, minimock.Diff(*mm_want_ptrs.message, mm_got.message))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSendMessage.t.Errorf("ProducerServiceMock.SendMessage got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSendMessage.SendMessageMock.defaultExpectation.results
		if mm_results == nil {
			mmSendMessage.t.Fatal("No results are set for the ProducerServiceMock.SendMessage")
		}
		return (*mm_results).err
	}
	if mmSendMessage.funcSendMessage != nil {
		return mmSendMessage.funcSendMessage(ctx, message)
	}
	mmSendMessage.t.Fatalf("Unexpected call to ProducerServiceMock.SendMessage. %v %v", ctx, message)
	return
}

// SendMessageAfterCounter returns a count of finished ProducerServiceMock.SendMessage invocations
func (mmSendMessage *ProducerServiceMock) SendMessageAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSendMessage.afterSendMessageCounter)
}

// SendMessageBeforeCounter returns a count of ProducerServiceMock.SendMessage invocations
func (mmSendMessage *ProducerServiceMock) SendMessageBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSendMessage.beforeSendMessageCounter)
}

// Calls returns a list of arguments used in each call to ProducerServiceMock.SendMessage.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSendMessage *mProducerServiceMockSendMessage) Calls() []*ProducerServiceMockSendMessageParams {
	mmSendMessage.mutex.RLock()

	argCopy := make([]*ProducerServiceMockSendMessageParams, len(mmSendMessage.callArgs))
	copy(argCopy, mmSendMessage.callArgs)

	mmSendMessage.mutex.RUnlock()

	return argCopy
}

// MinimockSendMessageDone returns true if the count of the SendMessage invocations corresponds
// the number of defined expectations
func (m *ProducerServiceMock) MinimockSendMessageDone() bool {
	if m.SendMessageMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.SendMessageMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.SendMessageMock.invocationsDone()
}

// MinimockSendMessageInspect logs each unmet expectation
func (m *ProducerServiceMock) MinimockSendMessageInspect() {
	for _, e := range m.SendMessageMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ProducerServiceMock.SendMessage with params: %#v", *e.params)
		}
	}

	afterSendMessageCounter := mm_atomic.LoadUint64(&m.afterSendMessageCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.SendMessageMock.defaultExpectation != nil && afterSendMessageCounter < 1 {
		if m.SendMessageMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ProducerServiceMock.SendMessage")
		} else {
			m.t.Errorf("Expected call to ProducerServiceMock.SendMessage with params: %#v", *m.SendMessageMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSendMessage != nil && afterSendMessageCounter < 1 {
		m.t.Error("Expected call to ProducerServiceMock.SendMessage")
	}

	if !m.SendMessageMock.invocationsDone() && afterSendMessageCounter > 0 {
		m.t.Errorf("Expected %d calls to ProducerServiceMock.SendMessage but found %d calls",
			mm_atomic.LoadUint64(&m.SendMessageMock.expectedInvocations), afterSendMessageCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ProducerServiceMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockSendMessageInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ProducerServiceMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ProducerServiceMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockSendMessageDone()
}
