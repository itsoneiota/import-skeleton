package sftp

import (
	"fmt"
	"testing"

	skeleton "github.com/itsoneiota/import-skeleton"
	"github.com/itsoneiota/metrics"
	"github.com/itsoneiota/ssftp-go"
)

func TestCanFindFile(t *testing.T) {
	c, err := ssftp.NewClientWithCredentials("localhost:9000", "elucid", "123")
	if err != nil {
		panic(err)
	}
	mmc := metrics.NewMockMetricsClient()
	mp := metrics.NewMetricPublisher(mmc)
	i := NewImporter(c, "/sftpdata/importStuff", mp)
	i.Poll(work)
}

func work(i skeleton.WorkItem) error {
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
