package test

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

const str = `dddddddddddddddddddddddddddddddddddddddddddddddd
sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss
qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq`

func TestSmallBufferSize(t *testing.T) {
	strReader := strings.NewReader(str)
	reader := bufio.NewReaderSize(strReader, 16)
	result, err := reader.ReadString('\n')
	t.Logf("result : %s \n", result)
	for err != io.EOF {
		result, err = reader.ReadString('\n')
		t.Logf("result : %s \n", result)
	}
}
