package utils

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_getPrecision(t *testing.T) {
	exam1 := "4999999999999999999999"
	exam2 := "49999999999999999999.99"
	exam3 := "49999999.999999.99999"
	exam4 := "49999999.999999"
	exam5 := "4999999999999.9999999"

	data := getPrecision(exam1)
	if data == "4999999999999990000000" {
		t.Log("OK")
	} else {
		t.Fatal("Failed")
	}
	data = getPrecision(exam2)
	if data == "49999999999999900000" {
		t.Log("OK")

	} else {
		t.Fatal("Failed")
	}
	data = getPrecision(exam3)
	if data == "49999999.999999.99999" {
		t.Log("OK")
	} else {
		t.Fatal("Failed")
	}
	data = getPrecision(exam4)
	if data == "49999999.999999" {
		t.Log("OK")
	} else {
		t.Fatal("Failed")
	}
	data = getPrecision(exam5)
	fmt.Println(data)
	amt, _ := strconv.ParseFloat(data, 64)
	fmt.Println(amt)
	if data == "4999999999999.999" {
		t.Log("OK")
	} else {
		t.Fatal("Failed")
	}
}

func TestMd5Encrypt(t *testing.T) {
	str := `{\"m_sequence\":\"16\",\"m_timeout\":\"16320\",\"m_source_port\":\"port-to-bank\",\"m_source_channel\":\"chann-to-gaia\",\"m_destination_port\":\"port-to-bank\",\"m_destination_channel\":\"chann-to-iris\",\"m_data\":\"eyJ0eXBlIjoiaWJjbW9ja2JhbmsvVHJhbnNmZXJQYWNrZXREYXRhIiwidmFsdWUiOnsiZGVub21pbmF0aW9uIjoidWlyaXMiLCJhbW91bnQiOiIxMDAwMDAwIiwic2VuZGVyIjoiZmFhMWVxdmtmdGh0cnI5M2c0cDlxc3BwNTR3NmR0anRybjI3OXZjbXBuIiwicmVjZWl2ZXIiOiJjb3Ntb3MxdDBsaHA5Mm1rZDhwZWpkcWp6bHo0bGN3eG1zM3U1Mzl1Y24zbHkiLCJzb3VyY2UiOnRydWV9fQ==\"}`
	data := []byte(str)
	res := Md5Encrypt(data)

	t.Log(res)
}
