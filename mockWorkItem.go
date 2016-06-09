package importer

const (
	fresh      = "FRESH"
	failed     = "FAILED"
	processing = "PROCESSING"
	completed  = "COMPLETED"
	terminated = "TERMINATED"
)

// MockWorkItem is a simple work item for use in tests.
type MockWorkItem struct {
	content string
	Status  string
}

// NewMockWorkItem returns a new mock work item with the given content.
func NewMockWorkItem(content string) *MockWorkItem {
	return &MockWorkItem{
		content: content,
		Status:  fresh,
	}
}

// Content gets the content of the work item.
func (m *MockWorkItem) Content() (string, error) {
	return m.content, nil
}

// Start indicates that work is starting on the item.
// Implementers should make sure that other workers cannot access the item.
func (m *MockWorkItem) Start() {
	m.Status = processing
}

// Complete marks a work item as complete.
// Further processing must be prevented.
// The item may be retained for archival purposes.
func (m *MockWorkItem) Complete(msg string) {
	m.Status = completed
}

// Fail indicates that work has not been successful, but the error may be recoverable.
// The item should be retained for retries, and the number of retries should be recorded if necessary.
func (m *MockWorkItem) Fail(msg string) {
	m.Status = failed
}

// Terminate indicates that work has not been successful, and that no further work can be done with the item.
// Further processing must be prevented.
// The item should be retained for review.
func (m *MockWorkItem) Terminate(msg string) {
	m.Status = terminated
}
