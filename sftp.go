package importer

import (
	"fmt"

	"github.com/itsoneiota/metrics"
	ssftp "github.com/itsoneiota/ssftp-go"
)

const (
	incoming   = "incoming"   // Directory for unprocessed files.
	processing = "processing" // Directory for files currently being processed.
	completed  = "completed"  // Directory for successfully processed files.
	terminated = "terminated" // Directory for files with unrecoverable errors.
)

// SFTPImporter represents an SFTP importer.
// We favour convention over configuration. An importer is rooted at a directory, which must have the following structure:
//
//     .
//     ├── completed   // Successfully completed files.
//     ├── incoming    // New, unprocessed files.
//     ├── processing  // Files being processed.
//     └── terminated  // Files with unrecoverable failures.
type SFTPImporter struct {
	client     *ssftp.Client
	worker     Worker
	metrics    *metrics.MetricPublisher
	incoming   string
	processing string
	completed  string
	terminated string
}

// Poll finds files ready for import.
func (i *SFTPImporter) Poll(w Worker) {
	// TODO: Mutex access to this function.
	items := i.findIncoming()
	fmt.Printf("Items To Import: %d", len(items))
	for _, item := range items {
		w(item)
	}
}

func (i *SFTPImporter) findIncoming() []WorkItem {
	w := i.client.Walk(i.incoming)
	var items []WorkItem
	for w.Step() {
		err := w.Err()
		if err != nil {
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

func (i *SFTPImporter) moveToProcessing(f *SFTPFile) {
	i.moveTo(f, i.processing)
}

func (i *SFTPImporter) moveToCompleted(f *SFTPFile) {
	i.moveTo(f, i.completed)
}

func (i *SFTPImporter) moveToTerminated(f *SFTPFile) {
	i.moveTo(f, i.terminated)
}

func (i *SFTPImporter) moveTo(f *SFTPFile, dst string) {
	newPath := dst + "/" + f.name
	i.client.Rename(f.path, newPath)
	f.path = newPath
}

// Content gets the string content of the file.
func (i *SFTPImporter) content(f *SFTPFile) (string, error) {
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

func (i *SFTPImporter) newFile(path string, name string) (*SFTPFile, error) {
	return &SFTPFile{
		importer: i,
		path:     path,
		name:     name,
	}, nil
}

// NewImporter returns a new importer using the given SFTP Client.
func NewImporter(c *ssftp.Client, dir string, m *metrics.MetricPublisher) Importer {
	return &SFTPImporter{
		client:     c,
		metrics:    m,
		incoming:   fmt.Sprintf("%s/%s", dir, incoming),
		processing: fmt.Sprintf("%s/%s", dir, processing),
		completed:  fmt.Sprintf("%s/%s", dir, completed),
		terminated: fmt.Sprintf("%s/%s", dir, terminated),
	}
}

// SFTPFile represents a file in an SFTP location.
type SFTPFile struct {
	importer *SFTPImporter
	path     string
	name     string
}

// Content gets the string content of the file.
func (f *SFTPFile) Content() (string, error) {
	return f.importer.content(f)
}

// Start moves a file to the 'processing' directory.
func (f *SFTPFile) Start() {
	f.importer.moveToProcessing(f)
}

// Complete moves a file to the 'archive' directory.
func (f *SFTPFile) Complete(msg string) {
	f.importer.moveToCompleted(f)
}

// Fail records a failure, for a later retry attempt.
func (f *SFTPFile) Fail(msg string) {
	// TODO: What do we do here? Move back to incoming?
}

// Terminate moves the file to the 'terminal' directory.
func (f *SFTPFile) Terminate(msg string) {
	f.importer.moveToTerminated(f)
}
