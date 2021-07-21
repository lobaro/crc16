package crc16

import (
	"testing"
)

type testCase struct {
	Message []byte
	CRC     uint16
}

func TestModbus(t *testing.T) {
	tests := []testCase{
		{[]byte{0xEA, 0x03, 0x00, 0x00, 0x00, 0x64}, 0x3A53},
		{[]byte{0x4B, 0x03, 0x00, 0x2C, 0x00, 0x37}, 0xBFCB},
		{[]byte{0x0D, 0x01, 0x00, 0x62, 0x00, 0x33}, 0x0DDD}}
	for _, testcase := range tests {
		result := ^ChecksumIBM(testcase.Message)
		if testcase.CRC != result {
			t.Fatalf("Modbus CRC-16 value is incorrect, expected %d, received %d.", testcase.CRC, result)
		}
	}
}

func BenchmarkChecksumIBM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ChecksumIBM([]byte{0xEA, 0x03, 0x00, 0x00, 0x00, 0x64})
	}
}

func BenchmarkMakeTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeTable(0xA001)
	}
}

func TestCCITTFalse(t *testing.T) {
	data := []byte("testdata")
	target := uint16(0xDC7C)

	actual := ChecksumCCITTFalse(data)
	if actual != target {
		t.Fatalf("CCITT checksum did not return the correct value, expected %x, received %x", target, actual)
	}
}

func TestXModem(t *testing.T) {
	tests := []testCase{
		{[]byte("123456789"), 0x31C3},
		{[]byte{0x54, 0x45, 0x4b, 0x38, 0x31, 0x31, 0x2c, 0x52, 0x35, 0x3d, 0x30, 0x31, 0x2c, 0x5a}, 0x8E48},
		{[]byte{0x54, 0x45, 0x4b, 0x38, 0x31, 0x31, 0x2c, 0x52, 0x35, 0x3d, 0x31, 0x31, 0x2c, 0x5a}, 0xF8FC}}
	for _, testcase := range tests {
		result := ChecksumXModem(testcase.Message)
		if testcase.CRC != result {
			t.Fatalf("Modbus CRC-16 value is incorrect, expected %d, received %d.", testcase.CRC, result)
		}
	}
}
