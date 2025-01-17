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

func hex2bin(hex string) []byte {
	bin := make([]byte, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		bin[i/2] = (byte(hex[i])-'0')<<4 + (byte(hex[i+1]) - '0')
	}
	return bin
}

func TestIBM_SDLC(t *testing.T) {
	tests := []testCase{
		{hex2bin("7602dd620062007263070177010a01a81502728412010201017777070100010800ff6401018201621e52fe5600008aca5a0177070100020800ff6401018201621e52fe5600000297d80177070100010801ff0101621e52fe5600008aca5a0177070100020801ff0101621e52fe5600000297d80177070100010802ff0101621e52fe5600000000000177070100020802ff0101621e52fe5600000000000177070100100700ff0101621b52fe5500000bc2010101"), 0xcb6f},
	}
	for _, testcase := range tests {
		result := ChecksumIBM_SDLC(testcase.Message)
		if testcase.CRC != result {
			t.Fatalf("IBM SDLC CRC-16 value is incorrect, expected %d, received %d.", testcase.CRC, result)
		}
	}
}
