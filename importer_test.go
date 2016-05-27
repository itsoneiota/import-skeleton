package importer

import "testing"

type dummyImporter int

func (dummyImporter) Poll(w Worker) {
	wi := dummyWorkItem(1)
	w(wi)
	return
}

type dummyWorkItem int

// Content
func (dummyWorkItem) Content() (string, error) {
	return "foo", nil
}

// Complete marks a work item as complete.
func (dummyWorkItem) Complete(msg string) {

}

// Start marks a work item as started.
func (dummyWorkItem) Start() {

}

// Fail indicates that no further work can be done on a work item.
func (dummyWorkItem) Fail(msg string) {

}

// Terminate indicates that no further work can be done on a work item.
func (dummyWorkItem) Terminate(msg string) {

}

var workerCalled bool

func dummyWorker(WorkItem) error {
	workerCalled = true
	return nil
}

func TestCanPoll(t *testing.T) {
	imp := dummyImporter(1)
	imp.Poll(dummyWorker)

	if !workerCalled {
		t.Errorf("Should have called the worker. Didn't.")
	}
}
