// Code generated by MockGen. DO NOT EDIT.
// Source: stevejefferson.co.uk/trac2gitea/accessor/trac (interfaces: Accessor)

// Package mock_trac is a generated GoMock package.
package mock_trac

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAccessor is a mock of Accessor interface
type MockAccessor struct {
	ctrl     *gomock.Controller
	recorder *MockAccessorMockRecorder
}

// MockAccessorMockRecorder is the mock recorder for MockAccessor
type MockAccessorMockRecorder struct {
	mock *MockAccessor
}

// NewMockAccessor creates a new mock instance
func NewMockAccessor(ctrl *gomock.Controller) *MockAccessor {
	mock := &MockAccessor{ctrl: ctrl}
	mock.recorder = &MockAccessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccessor) EXPECT() *MockAccessorMockRecorder {
	return m.recorder
}

// GetAttachments mocks base method
func (m *MockAccessor) GetAttachments(arg0 int64, arg1 func(int64, int64, int64, string, string, string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttachments", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAttachments indicates an expected call of GetAttachments
func (mr *MockAccessorMockRecorder) GetAttachments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttachments", reflect.TypeOf((*MockAccessor)(nil).GetAttachments), arg0, arg1)
}

// GetCommentString mocks base method
func (m *MockAccessor) GetCommentString(arg0, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentString", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentString indicates an expected call of GetCommentString
func (mr *MockAccessorMockRecorder) GetCommentString(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentString", reflect.TypeOf((*MockAccessor)(nil).GetCommentString), arg0, arg1)
}

// GetComments mocks base method
func (m *MockAccessor) GetComments(arg0 int64, arg1 func(int64, int64, string, string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComments", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetComments indicates an expected call of GetComments
func (mr *MockAccessorMockRecorder) GetComments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComments", reflect.TypeOf((*MockAccessor)(nil).GetComments), arg0, arg1)
}

// GetComponentNames mocks base method
func (m *MockAccessor) GetComponentNames(arg0 func(string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComponentNames", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetComponentNames indicates an expected call of GetComponentNames
func (mr *MockAccessorMockRecorder) GetComponentNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComponentNames", reflect.TypeOf((*MockAccessor)(nil).GetComponentNames), arg0)
}

// GetFullPath mocks base method
func (m *MockAccessor) GetFullPath(arg0 ...string) string {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFullPath", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFullPath indicates an expected call of GetFullPath
func (mr *MockAccessorMockRecorder) GetFullPath(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFullPath", reflect.TypeOf((*MockAccessor)(nil).GetFullPath), arg0...)
}

// GetMilestones mocks base method
func (m *MockAccessor) GetMilestones(arg0 func(string, string, int64, int64) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMilestones", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetMilestones indicates an expected call of GetMilestones
func (mr *MockAccessorMockRecorder) GetMilestones(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMilestones", reflect.TypeOf((*MockAccessor)(nil).GetMilestones), arg0)
}

// GetPriorityNames mocks base method
func (m *MockAccessor) GetPriorityNames(arg0 func(string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriorityNames", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetPriorityNames indicates an expected call of GetPriorityNames
func (mr *MockAccessorMockRecorder) GetPriorityNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriorityNames", reflect.TypeOf((*MockAccessor)(nil).GetPriorityNames), arg0)
}

// GetResolutionNames mocks base method
func (m *MockAccessor) GetResolutionNames(arg0 func(string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResolutionNames", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetResolutionNames indicates an expected call of GetResolutionNames
func (mr *MockAccessorMockRecorder) GetResolutionNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResolutionNames", reflect.TypeOf((*MockAccessor)(nil).GetResolutionNames), arg0)
}

// GetSeverityNames mocks base method
func (m *MockAccessor) GetSeverityNames(arg0 func(string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeverityNames", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetSeverityNames indicates an expected call of GetSeverityNames
func (mr *MockAccessorMockRecorder) GetSeverityNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeverityNames", reflect.TypeOf((*MockAccessor)(nil).GetSeverityNames), arg0)
}

// GetStringConfig mocks base method
func (m *MockAccessor) GetStringConfig(arg0, arg1 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringConfig", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetStringConfig indicates an expected call of GetStringConfig
func (mr *MockAccessorMockRecorder) GetStringConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringConfig", reflect.TypeOf((*MockAccessor)(nil).GetStringConfig), arg0, arg1)
}

// GetTicketAttachmentPath mocks base method
func (m *MockAccessor) GetTicketAttachmentPath(arg0 int64, arg1 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicketAttachmentPath", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTicketAttachmentPath indicates an expected call of GetTicketAttachmentPath
func (mr *MockAccessorMockRecorder) GetTicketAttachmentPath(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicketAttachmentPath", reflect.TypeOf((*MockAccessor)(nil).GetTicketAttachmentPath), arg0, arg1)
}

// GetTickets mocks base method
func (m *MockAccessor) GetTickets(arg0 func(int64, string, int64, string, string, string, string, string, string, string, string, string, string, string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTickets", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetTickets indicates an expected call of GetTickets
func (mr *MockAccessorMockRecorder) GetTickets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTickets", reflect.TypeOf((*MockAccessor)(nil).GetTickets), arg0)
}

// GetTypeNames mocks base method
func (m *MockAccessor) GetTypeNames(arg0 func(string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTypeNames", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetTypeNames indicates an expected call of GetTypeNames
func (mr *MockAccessorMockRecorder) GetTypeNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTypeNames", reflect.TypeOf((*MockAccessor)(nil).GetTypeNames), arg0)
}

// GetVersionNames mocks base method
func (m *MockAccessor) GetVersionNames(arg0 func(string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersionNames", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetVersionNames indicates an expected call of GetVersionNames
func (mr *MockAccessorMockRecorder) GetVersionNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersionNames", reflect.TypeOf((*MockAccessor)(nil).GetVersionNames), arg0)
}

// GetWikiAttachmentPath mocks base method
func (m *MockAccessor) GetWikiAttachmentPath(arg0, arg1 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWikiAttachmentPath", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetWikiAttachmentPath indicates an expected call of GetWikiAttachmentPath
func (mr *MockAccessorMockRecorder) GetWikiAttachmentPath(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWikiAttachmentPath", reflect.TypeOf((*MockAccessor)(nil).GetWikiAttachmentPath), arg0, arg1)
}

// GetWikiAttachments mocks base method
func (m *MockAccessor) GetWikiAttachments(arg0 func(string, string) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWikiAttachments", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetWikiAttachments indicates an expected call of GetWikiAttachments
func (mr *MockAccessorMockRecorder) GetWikiAttachments(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWikiAttachments", reflect.TypeOf((*MockAccessor)(nil).GetWikiAttachments), arg0)
}

// GetWikiPages mocks base method
func (m *MockAccessor) GetWikiPages(arg0 func(string, string, string, string, int64, int64) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWikiPages", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetWikiPages indicates an expected call of GetWikiPages
func (mr *MockAccessorMockRecorder) GetWikiPages(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWikiPages", reflect.TypeOf((*MockAccessor)(nil).GetWikiPages), arg0)
}

// IsPredefinedPage mocks base method
func (m *MockAccessor) IsPredefinedPage(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPredefinedPage", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPredefinedPage indicates an expected call of IsPredefinedPage
func (mr *MockAccessorMockRecorder) IsPredefinedPage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPredefinedPage", reflect.TypeOf((*MockAccessor)(nil).IsPredefinedPage), arg0)
}
