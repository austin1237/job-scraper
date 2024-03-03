package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDB struct {
	mock.Mock
}

func (m *MockDynamoDB) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

func (m *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func TestNewTable(t *testing.T) {
	table, err := NewTable("test", "us-west-2")
	assert.NoError(t, err)
	assert.NotNil(t, table)
}

func TestReadItem(t *testing.T) {
	mockSvc := new(MockDynamoDB)
	table := &Table{Name: "test", svc: mockSvc}

	mockSvc.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{}, nil)

	_, err := table.ReadItem("Acme Corp")
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
}

func TestWriteItems(t *testing.T) {
	mockSvc := new(MockDynamoDB)
	table := &Table{Name: "test", svc: mockSvc}

	mockSvc.On("UpdateItem", mock.Anything).Return(&dynamodb.UpdateItemOutput{}, nil)

	table.WriteItems([]string{"Acme Corp", "Globex Corporation"})

	mockSvc.AssertExpectations(t)
}
