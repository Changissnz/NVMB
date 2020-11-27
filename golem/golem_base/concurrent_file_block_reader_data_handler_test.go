package golem 

import (
	"testing"
)

func Test_ConcurrentFileBlockReader_MergeCSVFilesToFile(t *testing.T) {
	x := []string{
		"./datos/node_data/node_0",
		"./datos/node_data/node_1",
		"./datos/node_data/node_2",
		"./datos/node_data/node_3",
		"./datos/node_data/node_4",
		"./datos/node_data/node_5",
		"./datos/node_data/node_6",
		"./datos/node_data/node_7",
		"./datos/node_data/node_8",
		"./datos/node_data/node_9"}
		
	outfi := "./datos/node_data/mnd"
	c := OneConcurrentFileBlockReader("") 
	c.MergeCSVFilesToFile(x, outfi)
} 

func Test_ConcurrentFileBlockReader__NodeDataRead(t *testing.T) {
	x := "./datos/node_data/node_1"
	c := OneConcurrentFileBlockReader(x)
	c.ReadPartition("full")
}