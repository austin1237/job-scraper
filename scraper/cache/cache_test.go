package cache

import (
	"scraper/job"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTable struct {
	mock.Mock
}

func (m *MockTable) ReadItem(company string) (string, error) {
	args := m.Called(company)
	return args.String(0), args.Error(1)
}

func (m *MockTable) WriteItems(companies []string) {
	m.Called(companies)
}

func TestFilterCachedCompanies(t *testing.T) {
	mockTable := new(MockTable)
	mockTable.On("ReadItem", "Acme Corp").Return("Acme Corp", nil)
	mockTable.On("ReadItem", "Globex Corporation").Return("", nil)

	cache := &Cache{
		table: mockTable,
	}

	// Test the FilterCachedCompanies method
	jobs := []job.Job{
		{Company: "Acme Corp"},
		{Company: "Globex Corporation"},
	}
	notInCache, err := cache.FilterCachedCompanies(jobs)

	assert.NoError(t, err)
	assert.Len(t, notInCache, 1)
	assert.Equal(t, "Globex Corporation", notInCache[0].Company)

	mockTable.AssertExpectations(t)
}

func TestWriteCompaniesToCache(t *testing.T) {
	mockTable := new(MockTable)
	mockTable.On("WriteItems", []string{"Acme Corp", "Globex Corporation"}).Return()

	cache := &Cache{
		table: mockTable,
	}

	// Test the WriteCompaniesToCache method
	jobs := []job.Job{
		{Company: "Acme Corp"},
		{Company: "Globex Corporation"},
	}
	cache.WriteCompaniesToCache(jobs)

	mockTable.AssertExpectations(t)
}
