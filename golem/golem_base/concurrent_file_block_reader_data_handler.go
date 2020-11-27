package golem

import (
	"os"
	"fmt"
)

/*
*/ 
func (c *ConcurrentFileBlockReader) MergeCSVFilesToFile(fps []string, outfi string) bool {
	// check if file exists
	var fi *os.File
	var err error
	fi, err = os.OpenFile(outfi, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		/*
		fmt.Printf("invalid file path @ %s\n", outfi)
		fi, err = os.OpenFile(outfi, os.O_CREATE, 0600) 
		fi, err = os.OpenFile(outfi, os.O_APPEND, 0777)
		*/
		return false
	}

	for _, f := range fps {
		c.WriteDataOneFile(f, fi) 
	}

	return true
}

/*
*/ 
func (c *ConcurrentFileBlockReader) WriteDataOneFile(fp string, fi *os.File) { 
	columnsBefore := c.columns
	c.fp = fp

	c.ReadHeader()
	
	if columnsBefore != nil {
		if !EqualStringSlices(columnsBefore, c.columns) {

			fmt.Println("columns before") 
			fmt.Println(columnsBefore) 

			fmt.Println("columns after")
			fmt.Println(c.columns) 

			panic("file has different header than previous") 
		}
	} else {
		// write header out to file
		x := [][]string{c.columns}
		CSVDataToFileObject(fi, x)
	}
	
	c.PrelimRead()
	sz := int64(0) 
	for sz < c.fileLength {
		sz += c.ReadPartition("full")
		// write each block to
		c.WriteBlockDataToFile(fi)
	}

}

func (c *ConcurrentFileBlockReader) WriteBlockDataToFile(fi *os.File) {
	for _, b := range c.blockData {
		CSVDataToFileObject(fi, b.datos)  
	}
}