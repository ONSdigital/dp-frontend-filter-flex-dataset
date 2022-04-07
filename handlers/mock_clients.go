// Code generated by MockGen. DO NOT EDIT.
// Source: clients.go

// Package handlers is a generated GoMock package.
package handlers

import (
	context "context"
	io "io"
	reflect "reflect"

	dataset "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	dimension "github.com/ONSdigital/dp-api-clients-go/v2/dimension"
	filter "github.com/ONSdigital/dp-api-clients-go/v2/filter"
	model "github.com/ONSdigital/dp-renderer/model"
	gomock "github.com/golang/mock/gomock"
)

// MockClientError is a mock of ClientError interface.
type MockClientError struct {
	ctrl     *gomock.Controller
	recorder *MockClientErrorMockRecorder
}

// MockClientErrorMockRecorder is the mock recorder for MockClientError.
type MockClientErrorMockRecorder struct {
	mock *MockClientError
}

// NewMockClientError creates a new mock instance.
func NewMockClientError(ctrl *gomock.Controller) *MockClientError {
	mock := &MockClientError{ctrl: ctrl}
	mock.recorder = &MockClientErrorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientError) EXPECT() *MockClientErrorMockRecorder {
	return m.recorder
}

// Code mocks base method.
func (m *MockClientError) Code() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Code")
	ret0, _ := ret[0].(int)
	return ret0
}

// Code indicates an expected call of Code.
func (mr *MockClientErrorMockRecorder) Code() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Code", reflect.TypeOf((*MockClientError)(nil).Code))
}

// Error mocks base method.
func (m *MockClientError) Error() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error")
	ret0, _ := ret[0].(string)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *MockClientErrorMockRecorder) Error() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockClientError)(nil).Error))
}

// MockRenderClient is a mock of RenderClient interface.
type MockRenderClient struct {
	ctrl     *gomock.Controller
	recorder *MockRenderClientMockRecorder
}

// MockRenderClientMockRecorder is the mock recorder for MockRenderClient.
type MockRenderClientMockRecorder struct {
	mock *MockRenderClient
}

// NewMockRenderClient creates a new mock instance.
func NewMockRenderClient(ctrl *gomock.Controller) *MockRenderClient {
	mock := &MockRenderClient{ctrl: ctrl}
	mock.recorder = &MockRenderClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRenderClient) EXPECT() *MockRenderClientMockRecorder {
	return m.recorder
}

// BuildPage mocks base method.
func (m *MockRenderClient) BuildPage(w io.Writer, pageModel interface{}, templateName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BuildPage", w, pageModel, templateName)
}

// BuildPage indicates an expected call of BuildPage.
func (mr *MockRenderClientMockRecorder) BuildPage(w, pageModel, templateName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildPage", reflect.TypeOf((*MockRenderClient)(nil).BuildPage), w, pageModel, templateName)
}

// NewBasePageModel mocks base method.
func (m *MockRenderClient) NewBasePageModel() model.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBasePageModel")
	ret0, _ := ret[0].(model.Page)
	return ret0
}

// NewBasePageModel indicates an expected call of NewBasePageModel.
func (mr *MockRenderClientMockRecorder) NewBasePageModel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBasePageModel", reflect.TypeOf((*MockRenderClient)(nil).NewBasePageModel))
}

// MockFilterClient is a mock of FilterClient interface.
type MockFilterClient struct {
	ctrl     *gomock.Controller
	recorder *MockFilterClientMockRecorder
}

// MockFilterClientMockRecorder is the mock recorder for MockFilterClient.
type MockFilterClientMockRecorder struct {
	mock *MockFilterClient
}

// NewMockFilterClient creates a new mock instance.
func NewMockFilterClient(ctrl *gomock.Controller) *MockFilterClient {
	mock := &MockFilterClient{ctrl: ctrl}
	mock.recorder = &MockFilterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilterClient) EXPECT() *MockFilterClientMockRecorder {
	return m.recorder
}

// GetJobState mocks base method.
func (m *MockFilterClient) GetJobState(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterID string) (filter.Model, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobState", ctx, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterID)
	ret0, _ := ret[0].(filter.Model)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetJobState indicates an expected call of GetJobState.
func (mr *MockFilterClientMockRecorder) GetJobState(ctx, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobState", reflect.TypeOf((*MockFilterClient)(nil).GetJobState), ctx, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterID)
}

// UpdateFlexBlueprint mocks base method.
func (m_2 *MockFilterClient) UpdateFlexBlueprint(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID string, m filter.Model, doSubmit bool, populationType, ifMatch string) (filter.Model, string, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "UpdateFlexBlueprint", ctx, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, m, doSubmit, populationType, ifMatch)
	ret0, _ := ret[0].(filter.Model)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateFlexBlueprint indicates an expected call of UpdateFlexBlueprint.
func (mr *MockFilterClientMockRecorder) UpdateFlexBlueprint(ctx, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, m, doSubmit, populationType, ifMatch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlexBlueprint", reflect.TypeOf((*MockFilterClient)(nil).UpdateFlexBlueprint), ctx, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, m, doSubmit, populationType, ifMatch)
}

// MockDatasetClient is a mock of DatasetClient interface.
type MockDatasetClient struct {
	ctrl     *gomock.Controller
	recorder *MockDatasetClientMockRecorder
}

// MockDatasetClientMockRecorder is the mock recorder for MockDatasetClient.
type MockDatasetClientMockRecorder struct {
	mock *MockDatasetClient
}

// NewMockDatasetClient creates a new mock instance.
func NewMockDatasetClient(ctrl *gomock.Controller) *MockDatasetClient {
	mock := &MockDatasetClient{ctrl: ctrl}
	mock.recorder = &MockDatasetClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatasetClient) EXPECT() *MockDatasetClientMockRecorder {
	return m.recorder
}

// GetOptions mocks base method.
func (m *MockDatasetClient) GetOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension string, q *dataset.QueryParams) (dataset.Options, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOptions", ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension, q)
	ret0, _ := ret[0].(dataset.Options)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOptions indicates an expected call of GetOptions.
func (mr *MockDatasetClientMockRecorder) GetOptions(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOptions", reflect.TypeOf((*MockDatasetClient)(nil).GetOptions), ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version, dimension, q)
}

// MockDimensionClient is a mock of DimensionClient interface.
type MockDimensionClient struct {
	ctrl     *gomock.Controller
	recorder *MockDimensionClientMockRecorder
}

// MockDimensionClientMockRecorder is the mock recorder for MockDimensionClient.
type MockDimensionClientMockRecorder struct {
	mock *MockDimensionClient
}

// NewMockDimensionClient creates a new mock instance.
func NewMockDimensionClient(ctrl *gomock.Controller) *MockDimensionClient {
	mock := &MockDimensionClient{ctrl: ctrl}
	mock.recorder = &MockDimensionClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDimensionClient) EXPECT() *MockDimensionClientMockRecorder {
	return m.recorder
}

// GetAreaTypes mocks base method.
func (m *MockDimensionClient) GetAreaTypes(ctx context.Context, userAuthToken, serviceAuthToken, datasetID string) (dimension.GetAreaTypesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAreaTypes", ctx, userAuthToken, serviceAuthToken, datasetID)
	ret0, _ := ret[0].(dimension.GetAreaTypesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAreaTypes indicates an expected call of GetAreaTypes.
func (mr *MockDimensionClientMockRecorder) GetAreaTypes(ctx, userAuthToken, serviceAuthToken, datasetID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAreaTypes", reflect.TypeOf((*MockDimensionClient)(nil).GetAreaTypes), ctx, userAuthToken, serviceAuthToken, datasetID)
}
