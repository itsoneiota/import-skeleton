package sftp

import (
	"fmt"

	skeleton "github.com/itsoneiota/import-skeleton"
	"github.com/itsoneiota/metrics"
	ssftp "github.com/itsoneiota/ssftp-go"
)

const (
	incoming   = "incoming"   // Directory for unprocessed files.
	processing = "processing" // Directory for files currently being processed.
	completed  = "completed"  // Directory for successfully processed files.
	terminated = "terminated" // Directory for files with unrecoverable errors.
)

// Importer represents an SFTP importer.
// We favour convention over configuration. An importer is rooted at a directory, which must have the following structure:
//
//     .
//     ├── completed   // Successfully completed files.
//     ├── incoming    // New, unprocessed files.
//     ├── processing  // Files being processed.
//     └── terminated  // Files with unrecoverable failures.
type Importer struct {
	client     *ssftp.Client
	worker     skeleton.Worker
	metrics    *metrics.MetricPublisher
	incoming   string
	processing string
	completed  string
	terminated string
}

// Poll finds files ready for import.
func (i *Importer) Poll(w skeleton.Worker) {
	// TODO: Mutex access to this function.
	items := i.findIncoming()
	fmt.Printf("Items To Import: %d\r\n", len(items))
	for _, item := range items {
		w(item)
	}
}

func (i *Importer) findIncoming() []skeleton.WorkItem {
	w := i.client.Walk(i.incoming)
	var items []skeleton.WorkItem
	for w.Step() {
		err := w.Err()
		if err != nil {
			fmt.Printf("err: %s, path: %s\r\n", err, w.Path())
			continue
		}
		if w.Stat().IsDir() {
			continue
		}
		item, err := i.newFile(w.Path(), w.Stat().Name())
		if err != nil {
			continue
		}
		items = append(items, item)
		i.metrics.Client.Inc("IncomingItems", 1)
	}
	return items
}

func (i *Importer) moveToProcessing(f *File) {
	i.moveTo(f, i.processing)
}

func (i *Importer) moveToCompleted(f *File) {
	i.moveTo(f, i.completed)
}

func (i *Importer) moveToTerminated(f *File) {
	i.moveTo(f, i.terminated)
}

func (i *Importer) moveTo(f *File, dst string) {
	newPath := dst + "/" + f.name
	i.client.Rename(f.path, newPath)
	f.path = newPath
}

// Content gets the string content of the file.
func (i *Importer) content(f *File) (string, error) {
	file, err := i.client.Open(f.path)
	if err != nil {
		return "", err
	}
	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	bytes := make([]byte, info.Size())

	file.Read(bytes)
	str := string(bytes)
	return str, nil
}

func (i *Importer) newFile(path string, name string) (*File, error) {
	return &File{
		importer: i,
		path:     path,
		name:     name,
	}, nil
}

// NewImporter returns a new importer using the given SFTP Client.
func NewImporter(c *ssftp.Client, dir string, m *metrics.MetricPublisher) skeleton.Importer {
	return &Importer{
		client:     c,
		metrics:    m,
		incoming:   fmt.Sprintf("%s/%s", dir, incoming),
		processing: fmt.Sprintf("%s/%s", dir, processing),
		completed:  fmt.Sprintf("%s/%s", dir, completed),
		terminated: fmt.Sprintf("%s/%s", dir, terminated),
	}
}

// File represents a file in an SFTP location.
type File struct {
	importer *Importer
	path     string
	name     string
}

// Content gets the string content of the file.
func (f *File) Content() (string, error) {
	return f.importer.content(f)
}

// Start moves a file to the 'processing' directory.
func (f *File) Start() {
	f.importer.moveToProcessing(f)
}

// Complete moves a file to the 'archive' directory.
func (f *File) Complete(msg string) {
	f.importer.moveToCompleted(f)
	f.importer.metrics.Client.Inc("FileComplete", 1)

}

// Fail records a failure, for a later retry attempt.
func (f *File) Fail(msg string) {
	f.importer.metrics.Client.Inc("FileFailure", 1)
	// TODO: What do we do here? Move back to incoming?
}

// Terminate moves the file to the 'terminal' directory.
func (f *File) Terminate(msg string) {
	f.importer.moveToTerminated(f)
	f.importer.metrics.Client.Inc("FileTerminal", 1)
}
