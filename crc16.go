package crc16

// Predefined polynomials.
const (
	// IBM is used by Bisync, Modbus, USB, ANSI X3.28, SIA DC-07, ...
	IBM        = 0xA001
	IBM_NORMAL = 0x8005

	// CCITT is used by X.25, V.41, HDLC FCS, XMODEM, Bluetooth, PACTOR, SD, ...
	// CCITT forward is 0x8408. Reverse is 0x1021.
	CCITT      = 0x8408
	CCITTFalse = 0x1021

	// SCSI is used by SCSI
	SCSI = 0xEDD1

	MBUS = 0x3D65

	XMODEM = 0x1021
)

// Table is a 256-word table representing the polynomial for efficient processing.
type Table struct {
	entries  [256]uint16
	reversed bool
}

func (t *Table) Entries() [256]uint16 {
	return t.entries
}

// IBMTable is the table for the IBM polynomial.
var IBMTable = makeTable(IBM)

// CCITTTable is the table for the CCITT polynomial.
var CCITTTable = makeTable(CCITT)

// CCITTFalseTable is the table for CCITT-FALSE.
var CCITTFalseTable = makeBitsReversedTable(CCITTFalse)

// SCSITable is the table for the SCSI polynomial.
var SCSITable = makeTable(SCSI)

var MBusTable = makeBitsReversedTable(MBUS)

var XModemTable = makeBitsReversedTable(XMODEM)

var IBM_SDLCTable = makeTable(CCITTFalse)

// MakeTable returns the Table constructed from the specified polynomial.
func MakeTable(poly uint16) *Table {
	return makeTable(poly)
}

// MakeTable returns the Table constructed from the specified polynomial.
func MakeBitsReversedTable(poly uint16) *Table {
	return makeBitsReversedTable(poly)
}

// makeTable returns the Table constructed from the specified polynomial.
func makeBitsReversedTable(poly uint16) *Table {
	t := &Table{
		reversed: true,
	}
	width := uint16(16)
	for i := uint16(0); i < 256; i++ {
		crc := i << (width - 8)
		for j := 0; j < 8; j++ {
			if crc&(1<<(width-1)) != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		t.entries[i] = crc
	}
	return t
}

func makeTable(poly uint16) *Table {
	t := &Table{
		reversed: false,
	}
	for i := 0; i < 256; i++ {
		crc := uint16(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
		t.entries[i] = crc
	}
	return t
}

func updateBitsReversed(crc uint16, tab *Table, p []byte) uint16 {
	for _, v := range p {
		crc = tab.entries[byte(crc>>8)^v] ^ (crc << 8)
	}
	return crc
}

func update(crc uint16, tab *Table, p []byte) uint16 {
	//crc = ^crc

	for _, v := range p {
		crc = tab.entries[byte(crc)^v] ^ (crc >> 8)
	}

	return crc
}

// Update returns the result of adding the bytes in p to the crc.
func Update(crc uint16, tab *Table, p []byte) uint16 {
	if tab.reversed {
		return updateBitsReversed(crc, tab, p)
	} else {
		return update(crc, tab, p)
	}
}

func UpdateInverted(crc uint16, tab *Table, p []byte) uint16 {
	if tab.reversed {
		return ^updateBitsReversed(^crc, tab, p)
	} else {
		return ^update(^crc, tab, p)
	}
}

// Checksum returns the CRC-16 checksum of data
// using the polynomial represented by the Table.
func Checksum(data []byte, tab *Table) uint16 { return Update(0, tab, data) }

// ChecksumInverted returns the CRC-16 checksum of data
// where each byte gets inverted before lookup
// using the polynomial represented by the Table.
func ChecksumInverted(data []byte, tab *Table) uint16 { return UpdateInverted(0, tab, data) }

// ChecksumIBM returns the CRC-16 checksum of data
// using the IBM polynomial.
func ChecksumIBM(data []byte) uint16 { return UpdateInverted(0, IBMTable, data) }

// ChecksumCCITTFalse returns the CRC-16 checksum using
// what some call the CCITT-False polynomial, which matches what is used
// by Perl Digest/CRC and Boost for example.
func ChecksumCCITTFalse(data []byte) uint16 { return Update(0xffff, CCITTFalseTable, data) }

// ChecksumCCITT returns the CRC-16 checksum of data
// using the CCITT polynomial.
func ChecksumCCITT(data []byte) uint16 { return UpdateInverted(0, CCITTTable, data) }

// ChecksumSCSI returns the CRC-16 checksum of data
// using the SCSI polynomial.
func ChecksumSCSI(data []byte) uint16 { return UpdateInverted(0, SCSITable, data) }

// ChecksumMBus returns the CRC-16 checksum of data
// using the MBus polynomial. Final result is inverted
func ChecksumMBus(data []byte) uint16 { return ^Update(0, MBusTable, data) }

// ChecksumXModem returns the CRC-16 checksum of data
// using the XMODEM polynomial.
func ChecksumXModem(data []byte) uint16 { return Update(0, XModemTable, data) }

// ChecksumIBM_SDLC returns the CRC-16 checksum of data
// using the IBM-SDLC polynomial.
func ChecksumIBM_SDLC(data []byte) uint16 { return Update(0xffff, IBM_SDLCTable, data) }
