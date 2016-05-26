package importer

import (
	"fmt"

	ssftp "github.com/itsoneiota/ssftp-go"
	"github.com/pkg/sftp"
)

const (
	incoming   = "incoming"   // Directory for unprocessed files.
	processing = "processing" // Directory for files currently being processed.
	completed  = "completed"  // Directory for successfully processed files.
	terminated = "terminated" // Directory for files with unrecoverable errors.
)

// SFTPImporter represents an SFTP importer.
type SFTPImporter struct {
	client     *ssftp.Client
	worker     Worker
	incoming   string
	processing string
	completed  string
	terminated string
}

// Poll finds files ready for import.
func (i *SFTPImporter) Poll(w Worker) {
	// TODO: Mutex access to this function.
	items := i.findIncoming()
	for _, item := range items {
		w.Handle(item)
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
		item, err := i.newSFTPFile(w.Path())
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items
}

// NewSFTPFile returns an SFTPFile
func (i *SFTPImporter) newSFTPFile(path string) (*SFTPFile, error) {
	file, err := i.client.Open(path)
	if err != nil {
		return nil, err
	}
	return &SFTPFile{
		client: i.client,
		path:   path,
		file:   file,
	}, nil
}

// NewImporter returns a new importer using the given SFTP Client.
func NewImporter(c *ssftp.Client, dir string) Importer {
	return &SFTPImporter{
		client:     c,
		incoming:   fmt.Sprintf("%s/%s", dir, incoming),
		processing: fmt.Sprintf("%s/%s", dir, processing),
		completed:  fmt.Sprintf("%s/%s", dir, completed),
		terminated: fmt.Sprintf("%s/%s", dir, terminated),
	}
}

// SFTPFile represents a file in an SFTP location.
type SFTPFile struct {
	client *ssftp.Client
	file   *sftp.File
	path   string
}

// Content gets the string content of the file.
func (f *SFTPFile) Content() (string, error) {
	info, err := f.file.Stat()
	if err != nil {
		return "", err
	}
	
	bytes := make([]byte, info.Size())
	// info.Size()
	f.file.Read(bytes)
	str := string(bytes)
	return str, nil
}

// Start moves a file to the 'processing' directory.
func (f *SFTPFile) Start() {
	// TODO: Move file to processing.
}

// Complete moves a file to the 'archive' directory.
func (f *SFTPFile) Complete(msg string) {
	// TODO: Move file to completed.
}

// Fail records a failure, for a later retry attempt.
func (f *SFTPFile) Fail(msg string) {
	// TODO: Move file to failed location.
}

// Terminate moves the file to the 'terminal' directory.
func (f *SFTPFile) Terminate(msg string) {
	// TODO:
}
