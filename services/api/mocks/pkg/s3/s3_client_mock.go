// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_s3

import (
	context "context"

	s3 "github.com/Zampfi/application-platform/services/api/pkg/s3"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockS3Client is an autogenerated mock type for the S3Client type
type MockS3Client struct {
	mock.Mock
}

type MockS3Client_Expecter struct {
	mock *mock.Mock
}

func (_m *MockS3Client) EXPECT() *MockS3Client_Expecter {
	return &MockS3Client_Expecter{mock: &_m.Mock}
}

// CopyFile provides a mock function with given fields: ctx, srcBucket, srcKey, destBucket, destKey
func (_m *MockS3Client) CopyFile(ctx context.Context, srcBucket string, srcKey string, destBucket string, destKey string) error {
	ret := _m.Called(ctx, srcBucket, srcKey, destBucket, destKey)

	if len(ret) == 0 {
		panic("no return value specified for CopyFile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) error); ok {
		r0 = rf(ctx, srcBucket, srcKey, destBucket, destKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockS3Client_CopyFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CopyFile'
type MockS3Client_CopyFile_Call struct {
	*mock.Call
}

// CopyFile is a helper method to define mock.On call
//   - ctx context.Context
//   - srcBucket string
//   - srcKey string
//   - destBucket string
//   - destKey string
func (_e *MockS3Client_Expecter) CopyFile(ctx interface{}, srcBucket interface{}, srcKey interface{}, destBucket interface{}, destKey interface{}) *MockS3Client_CopyFile_Call {
	return &MockS3Client_CopyFile_Call{Call: _e.mock.On("CopyFile", ctx, srcBucket, srcKey, destBucket, destKey)}
}

func (_c *MockS3Client_CopyFile_Call) Run(run func(ctx context.Context, srcBucket string, srcKey string, destBucket string, destKey string)) *MockS3Client_CopyFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].(string))
	})
	return _c
}

func (_c *MockS3Client_CopyFile_Call) Return(_a0 error) *MockS3Client_CopyFile_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockS3Client_CopyFile_Call) RunAndReturn(run func(context.Context, string, string, string, string) error) *MockS3Client_CopyFile_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateUploadURL provides a mock function with given fields: ctx, bucket, key, expiry, contentType
func (_m *MockS3Client) GenerateUploadURL(ctx context.Context, bucket string, key string, expiry time.Time, contentType string) (string, error) {
	ret := _m.Called(ctx, bucket, key, expiry, contentType)

	if len(ret) == 0 {
		panic("no return value specified for GenerateUploadURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time, string) (string, error)); ok {
		return rf(ctx, bucket, key, expiry, contentType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time, string) string); ok {
		r0 = rf(ctx, bucket, key, expiry, contentType)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, time.Time, string) error); ok {
		r1 = rf(ctx, bucket, key, expiry, contentType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockS3Client_GenerateUploadURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateUploadURL'
type MockS3Client_GenerateUploadURL_Call struct {
	*mock.Call
}

// GenerateUploadURL is a helper method to define mock.On call
//   - ctx context.Context
//   - bucket string
//   - key string
//   - expiry time.Time
//   - contentType string
func (_e *MockS3Client_Expecter) GenerateUploadURL(ctx interface{}, bucket interface{}, key interface{}, expiry interface{}, contentType interface{}) *MockS3Client_GenerateUploadURL_Call {
	return &MockS3Client_GenerateUploadURL_Call{Call: _e.mock.On("GenerateUploadURL", ctx, bucket, key, expiry, contentType)}
}

func (_c *MockS3Client_GenerateUploadURL_Call) Run(run func(ctx context.Context, bucket string, key string, expiry time.Time, contentType string)) *MockS3Client_GenerateUploadURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(time.Time), args[4].(string))
	})
	return _c
}

func (_c *MockS3Client_GenerateUploadURL_Call) Return(_a0 string, _a1 error) *MockS3Client_GenerateUploadURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockS3Client_GenerateUploadURL_Call) RunAndReturn(run func(context.Context, string, string, time.Time, string) (string, error)) *MockS3Client_GenerateUploadURL_Call {
	_c.Call.Return(run)
	return _c
}

// GetFileDetails provides a mock function with given fields: ctx, bucket, key
func (_m *MockS3Client) GetFileDetails(ctx context.Context, bucket string, key string) (*s3.FileInfo, error) {
	ret := _m.Called(ctx, bucket, key)

	if len(ret) == 0 {
		panic("no return value specified for GetFileDetails")
	}

	var r0 *s3.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*s3.FileInfo, error)); ok {
		return rf(ctx, bucket, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *s3.FileInfo); ok {
		r0 = rf(ctx, bucket, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, bucket, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockS3Client_GetFileDetails_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFileDetails'
type MockS3Client_GetFileDetails_Call struct {
	*mock.Call
}

// GetFileDetails is a helper method to define mock.On call
//   - ctx context.Context
//   - bucket string
//   - key string
func (_e *MockS3Client_Expecter) GetFileDetails(ctx interface{}, bucket interface{}, key interface{}) *MockS3Client_GetFileDetails_Call {
	return &MockS3Client_GetFileDetails_Call{Call: _e.mock.On("GetFileDetails", ctx, bucket, key)}
}

func (_c *MockS3Client_GetFileDetails_Call) Run(run func(ctx context.Context, bucket string, key string)) *MockS3Client_GetFileDetails_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockS3Client_GetFileDetails_Call) Return(_a0 *s3.FileInfo, _a1 error) *MockS3Client_GetFileDetails_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockS3Client_GetFileDetails_Call) RunAndReturn(run func(context.Context, string, string) (*s3.FileInfo, error)) *MockS3Client_GetFileDetails_Call {
	_c.Call.Return(run)
	return _c
}

// GetSampleFilePathFromFolder provides a mock function with given fields: ctx, bucket, folderPrefix
func (_m *MockS3Client) GetSampleFilePathFromFolder(ctx context.Context, bucket string, folderPrefix string) (string, error) {
	ret := _m.Called(ctx, bucket, folderPrefix)

	if len(ret) == 0 {
		panic("no return value specified for GetSampleFilePathFromFolder")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, bucket, folderPrefix)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, bucket, folderPrefix)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, bucket, folderPrefix)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockS3Client_GetSampleFilePathFromFolder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSampleFilePathFromFolder'
type MockS3Client_GetSampleFilePathFromFolder_Call struct {
	*mock.Call
}

// GetSampleFilePathFromFolder is a helper method to define mock.On call
//   - ctx context.Context
//   - bucket string
//   - folderPrefix string
func (_e *MockS3Client_Expecter) GetSampleFilePathFromFolder(ctx interface{}, bucket interface{}, folderPrefix interface{}) *MockS3Client_GetSampleFilePathFromFolder_Call {
	return &MockS3Client_GetSampleFilePathFromFolder_Call{Call: _e.mock.On("GetSampleFilePathFromFolder", ctx, bucket, folderPrefix)}
}

func (_c *MockS3Client_GetSampleFilePathFromFolder_Call) Run(run func(ctx context.Context, bucket string, folderPrefix string)) *MockS3Client_GetSampleFilePathFromFolder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockS3Client_GetSampleFilePathFromFolder_Call) Return(_a0 string, _a1 error) *MockS3Client_GetSampleFilePathFromFolder_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockS3Client_GetSampleFilePathFromFolder_Call) RunAndReturn(run func(context.Context, string, string) (string, error)) *MockS3Client_GetSampleFilePathFromFolder_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockS3Client creates a new instance of MockS3Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockS3Client(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockS3Client {
	mock := &MockS3Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
