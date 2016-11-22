package cdb

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

var BadFormatError = errors.New("bad format")

// Make reads cdb-formatted records from r and writes a cdb-format database
// to w.  See the documentation for Dump for details on the input record format.
func Make(w io.WriteSeeker, r io.Reader) (err error) {
	defer func() { // Centralize error handling.
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	if _, err = w.Seek(int64(headerSize), 0); err != nil {
		return
	}

	buf := make([]byte, 8)
	rb := bufio.NewScanner(r)
	wb := bufio.NewWriter(w)
	hash := cdbHash()
	hw := io.MultiWriter(hash, wb) // Computes hash when writing record key.
	rr := &recReader{rb}

	htables := make(map[uint32][]slot)
	pos := headerSize

	// Read all records and write to output.
	for {
		kv, _, r := rr.readKV()

		if r == false {
			break
		}
		if kv == nil {
			continue
		}

		// Record format is "+klen,dlen:key->data\n"
		klen, dlen := uint32(len(kv[0])), uint32(len(kv[1]))

		writeNums(wb, klen, dlen, buf)
		hash.Reset()
		copyn(hw, kv[0])
		copyn(wb, kv[1])
		h := hash.Sum32()
		tableNum := h % 256
		htables[tableNum] = append(htables[tableNum], slot{h, pos})
		pos += 8 + klen + dlen
	}

	// Write hash tables and header.

	// Create and reuse a single hash table.
	maxSlots := 0
	for _, slots := range htables {
		if len(slots) > maxSlots {
			maxSlots = len(slots)
		}
	}
	slotTable := make([]slot, maxSlots*2)

	header := make([]byte, headerSize)
	// Write hash tables.
	for i := uint32(0); i < 256; i++ {
		slots := htables[i]
		if slots == nil {
			putNum(header[i*8:], pos)
			continue
		}

		nslots := uint32(len(slots) * 2)
		hashSlotTable := slotTable[:nslots]
		// Reset table slots.
		for j := 0; j < len(hashSlotTable); j++ {
			hashSlotTable[j].h = 0
			hashSlotTable[j].pos = 0
		}

		for _, slot := range slots {
			slotPos := (slot.h / 256) % nslots
			for hashSlotTable[slotPos].pos != 0 {
				slotPos++
				if slotPos == uint32(len(hashSlotTable)) {
					slotPos = 0
				}
			}
			hashSlotTable[slotPos] = slot
		}

		if err = writeSlots(wb, hashSlotTable, buf); err != nil {
			return
		}

		putNum(header[i*8:], pos)
		putNum(header[i*8+4:], nslots)
		pos += 8 * nslots
	}

	if err = wb.Flush(); err != nil {
		return
	}

	if _, err = w.Seek(0, 0); err != nil {
		return
	}

	_, err = w.Write(header)

	return
}

type recReader struct {
	*bufio.Scanner
}

// Return an array reading and decomposing the data from the scanner
// @todo: Change from string to []byte
func (rr *recReader) readKV() ([]string, string, bool) {

	if rr.Scan() == false {
		return nil, "", false
	}

	s := rr.Text()
	if len(s) == 0 || strings.HasPrefix(s, "#") == true {
		return nil, s, true
	}

	sarr := strings.Split(s, "=")
	sret := make([]string, len(sarr))
	var ctr int
	for _, v := range sarr {
		if v != "=" {
			sret[ctr] = strings.TrimSpace(v)
			ctr++
		}
	}

	if len(sret) != 2 {
		fmt.Println("We have a problem", sret, len(sret), len(s))
	}

	return sret, "", true
}

func copyn(w io.Writer, s string) {
	if _, err := io.WriteString(w, s); err != nil {
		panic(err)
	}
}

func putNum(buf []byte, x uint32) {
	binary.LittleEndian.PutUint32(buf, x)
}

func writeNums(w io.Writer, x, y uint32, buf []byte) {
	putNum(buf, x)
	putNum(buf[4:], y)
	if _, err := w.Write(buf[:8]); err != nil {
		panic(err)
	}
}

type slot struct {
	h, pos uint32
}

func writeSlots(w io.Writer, slots []slot, buf []byte) (err error) {
	for _, np := range slots {
		putNum(buf, np.h)
		putNum(buf[4:], np.pos)
		if _, err = w.Write(buf[:8]); err != nil {
			return
		}
	}

	return nil
}
