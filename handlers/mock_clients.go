// Code generated by MockGen. DO NOT EDIT.
// Source: clients.go

// Package handlers is a generated GoMock package.
package handlers

import (
	context "context"
	io "io"
	reflect "reflect"

	dataset "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	filter "github.com/ONSdigital/dp-api-clients-go/v2/filter"
	population "github.com/ONSdigital/dp-api-clients-go/v2/population"
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

// AddDimensionValue mocks base method.
func (m *MockFilterClient) AddDimensionValue(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDimensionValue", ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDimensionValue indicates an expected call of AddDimensionValue.
func (mr *MockFilterClientMockRecorder) AddDimensionValue(ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDimensionValue", reflect.TypeOf((*MockFilterClient)(nil).AddDimensionValue), ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch)
}

// DeleteDimensionOptions mocks base method.
func (m *MockFilterClient) DeleteDimensionOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDimensionOptions", ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteDimensionOptions indicates an expected call of DeleteDimensionOptions.
func (mr *MockFilterClientMockRecorder) DeleteDimensionOptions(ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDimensionOptions", reflect.TypeOf((*MockFilterClient)(nil).DeleteDimensionOptions), ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name)
}

// GetDimension mocks base method.
func (m *MockFilterClient) GetDimension(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name string) (filter.Dimension, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDimension", ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name)
	ret0, _ := ret[0].(filter.Dimension)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDimension indicates an expected call of GetDimension.
func (mr *MockFilterClientMockRecorder) GetDimension(ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDimension", reflect.TypeOf((*MockFilterClient)(nil).GetDimension), ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name)
}

// GetDimensionOptions mocks base method.
func (m *MockFilterClient) GetDimensionOptions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name string, q *filter.QueryParams) (filter.DimensionOptions, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDimensionOptions", ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, q)
	ret0, _ := ret[0].(filter.DimensionOptions)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDimensionOptions indicates an expected call of GetDimensionOptions.
func (mr *MockFilterClientMockRecorder) GetDimensionOptions(ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDimensionOptions", reflect.TypeOf((*MockFilterClient)(nil).GetDimensionOptions), ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, q)
}

// GetDimensions mocks base method.
func (m *MockFilterClient) GetDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID string, q *filter.QueryParams) (filter.Dimensions, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDimensions", ctx, userAuthToken, serviceAuthToken, collectionID, filterID, q)
	ret0, _ := ret[0].(filter.Dimensions)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetDimensions indicates an expected call of GetDimensions.
func (mr *MockFilterClientMockRecorder) GetDimensions(ctx, userAuthToken, serviceAuthToken, collectionID, filterID, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDimensions", reflect.TypeOf((*MockFilterClient)(nil).GetDimensions), ctx, userAuthToken, serviceAuthToken, collectionID, filterID, q)
}

// GetFilter mocks base method.
func (m *MockFilterClient) GetFilter(ctx context.Context, input filter.GetFilterInput) (*filter.GetFilterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilter", ctx, input)
	ret0, _ := ret[0].(*filter.GetFilterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilter indicates an expected call of GetFilter.
func (mr *MockFilterClientMockRecorder) GetFilter(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilter", reflect.TypeOf((*MockFilterClient)(nil).GetFilter), ctx, input)
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

// RemoveDimensionValue mocks base method.
func (m *MockFilterClient) RemoveDimensionValue(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDimensionValue", ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveDimensionValue indicates an expected call of RemoveDimensionValue.
func (mr *MockFilterClientMockRecorder) RemoveDimensionValue(ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDimensionValue", reflect.TypeOf((*MockFilterClient)(nil).RemoveDimensionValue), ctx, userAuthToken, serviceAuthToken, collectionID, filterID, name, value, ifMatch)
}

// SubmitFilter mocks base method.
func (m *MockFilterClient) SubmitFilter(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, ifMatch string, sfr filter.SubmitFilterRequest) (*filter.SubmitFilterResponse, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitFilter", ctx, userAuthToken, serviceAuthToken, downloadServiceToken, ifMatch, sfr)
	ret0, _ := ret[0].(*filter.SubmitFilterResponse)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SubmitFilter indicates an expected call of SubmitFilter.
func (mr *MockFilterClientMockRecorder) SubmitFilter(ctx, userAuthToken, serviceAuthToken, downloadServiceToken, ifMatch, sfr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitFilter", reflect.TypeOf((*MockFilterClient)(nil).SubmitFilter), ctx, userAuthToken, serviceAuthToken, downloadServiceToken, ifMatch, sfr)
}

// UpdateDimensions mocks base method.
func (m *MockFilterClient) UpdateDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, name, ifMatch string, dimension filter.Dimension) (filter.Dimension, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDimensions", ctx, userAuthToken, serviceAuthToken, collectionID, id, name, ifMatch, dimension)
	ret0, _ := ret[0].(filter.Dimension)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateDimensions indicates an expected call of UpdateDimensions.
func (mr *MockFilterClientMockRecorder) UpdateDimensions(ctx, userAuthToken, serviceAuthToken, collectionID, id, name, ifMatch, dimension interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDimensions", reflect.TypeOf((*MockFilterClient)(nil).UpdateDimensions), ctx, userAuthToken, serviceAuthToken, collectionID, id, name, ifMatch, dimension)
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

// GetVersionDimensions mocks base method.
func (m *MockDatasetClient) GetVersionDimensions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, id, edition, version string) (dataset.VersionDimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersionDimensions", ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version)
	ret0, _ := ret[0].(dataset.VersionDimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVersionDimensions indicates an expected call of GetVersionDimensions.
func (mr *MockDatasetClientMockRecorder) GetVersionDimensions(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersionDimensions", reflect.TypeOf((*MockDatasetClient)(nil).GetVersionDimensions), ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version)
}

// MockPopulationClient is a mock of PopulationClient interface.
type MockPopulationClient struct {
	ctrl     *gomock.Controller
	recorder *MockPopulationClientMockRecorder
}

// MockPopulationClientMockRecorder is the mock recorder for MockPopulationClient.
type MockPopulationClientMockRecorder struct {
	mock *MockPopulationClient
}

// NewMockPopulationClient creates a new mock instance.
func NewMockPopulationClient(ctrl *gomock.Controller) *MockPopulationClient {
	mock := &MockPopulationClient{ctrl: ctrl}
	mock.recorder = &MockPopulationClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPopulationClient) EXPECT() *MockPopulationClientMockRecorder {
	return m.recorder
}

// GetAreas mocks base method.
func (m *MockPopulationClient) GetAreas(ctx context.Context, input population.GetAreasInput) (population.GetAreasResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAreas", ctx, input)
	ret0, _ := ret[0].(population.GetAreasResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAreas indicates an expected call of GetAreas.
func (mr *MockPopulationClientMockRecorder) GetAreas(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAreas", reflect.TypeOf((*MockPopulationClient)(nil).GetAreas), ctx, input)
}

// GetPopulationAreaTypes mocks base method.
func (m *MockPopulationClient) GetPopulationAreaTypes(ctx context.Context, userAuthToken, serviceAuthToken, datasetID string) (population.GetAreaTypesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPopulationAreaTypes", ctx, userAuthToken, serviceAuthToken, datasetID)
	ret0, _ := ret[0].(population.GetAreaTypesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPopulationAreaTypes indicates an expected call of GetPopulationAreaTypes.
func (mr *MockPopulationClientMockRecorder) GetPopulationAreaTypes(ctx, userAuthToken, serviceAuthToken, datasetID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPopulationAreaTypes", reflect.TypeOf((*MockPopulationClient)(nil).GetPopulationAreaTypes), ctx, userAuthToken, serviceAuthToken, datasetID)
}
