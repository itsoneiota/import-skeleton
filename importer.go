package importer

// Importer describes a generic importer.
type Importer interface {
	// Polls for importable work.
	// If there is work to do, a WorkItem is passed to the Worker.
	Poll(Worker)
}

// Worker describes a function for handling the import of a single item.
type Worker func(WorkItem) error

// WorkItem represents an item of work for an importer.
// These may come from an (S)FTP directory, queue, or other transport.
type WorkItem interface {

	// Content gets the content of the work item.
	Content() (string, error)

	// Start indicates that work is starting on the item.
	// Implementers should make sure that other workers cannot access the item.
	Start()

	// Complete marks a work item as complete.
	// Further processing must be prevented.
	// The item may be retained for archival purposes.
	Complete(msg string)

	// Fail indicates that work has not been successful, but the error may be recoverable.
	// The item should be retained for retries, and the number of retries should be recorded if necessary.
	Fail(msg string)

	// Terminate indicates that work has not been successful, and that no further work can be done with the item.
	// Further processing must be prevented.
	// The item should be retained for review.
	Terminate(msg string)
}
