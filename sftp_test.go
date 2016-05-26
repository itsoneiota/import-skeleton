package importer

import (
	"fmt"

	"github.com/itsoneiota/ssftp-go"
	"testing"
)

func TestCanFindFile(t *testing.T) {
	c, err := ssftp.NewClientWithCredentials("localhost:9000", "elucid", "123")
	if err != nil {
		panic(err)
	}
	w := &echoWorker{}
	i := NewImporter(c, "/sftpdata/importStuff")
	i.Poll(w)
}

type echoWorker struct{}

func (*echoWorker) Handle(i WorkItem) {
	fmt.Println(i.Content())
}
