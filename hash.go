package crc16

import "hash"

// Hash16 is the common interface implemented by all 16-bit hash functions.
type Hash16 interface {
	hash.Hash
	Sum16() uint16
}

// New creates a new Hash16 computing the CRC-16 checksum
// using the polynomial represented by the Table.
func New(tab *Table) Hash16 { return &digest{0, tab} }

// digest represents the partial evaluation of a checksum.
type digest struct {
	crc uint16
	tab *Table
}

func (d *digest) Size() int { return 2 }

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Reset() { d.crc = 0 }

func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = Update(d.crc, d.tab, p)
	return len(p), nil
}

func (d *digest) Sum16() uint16 { return d.crc }

func (d *digest) Sum(in []byte) []byte {
	s := d.Sum16()
	return append(in, byte(s>>8), byte(s))
}
