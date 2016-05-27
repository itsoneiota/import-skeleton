package importer

import (
	"fmt"

	"testing"

	"github.com/itsoneiota/ssftp-go"
)

func TestCanFindFile(t *testing.T) {
	c, err := ssftp.NewClientWithCredentials("localhost:9000", "elucid", "123")
	if err != nil {
		panic(err)
	}
	i := NewImporter(c, "/sftpdata/importStuff")
	i.Poll(work)
}

func work(i WorkItem) error {
	fmt.Println("Here we go...")
	i.Start()
	content, _ := i.Content()
	fmt.Printf("workItem content: %s", content)
	if len(content) > 10 {
		fmt.Printf("Failing string %s", content)
		i.Terminate("Too long.")
	} else {
		fmt.Println("done")
		i.Complete("Completed successfully.")
	}
	return nil
}
