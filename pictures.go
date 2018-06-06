// pictures.go

// Copyright (C) 2018  Steve Merrony

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tello

// TakePicture requests the Tello to take a JPEG snapshot
func (tello *Tello) TakePicture() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoTakePic, tello.ctrlSeq, 0)
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

func (tello *Tello) sendFileSize() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	tello.ctrlSeq++
	tello.ctrlConn.Write(packetToBuffer(newPacket(ptData1, msgFileSize, tello.ctrlSeq, 1)))
}

func (tello *Tello) sendFileAckPiece(done byte, fID uint16, pieceNum uint32) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	tello.ctrlSeq++
	pkt := newPacket(ptData1, msgFileData, tello.ctrlSeq, 7)
	pkt.payload[0] = done
	pkt.payload[1] = byte(fID)
	pkt.payload[2] = byte(fID >> 8)
	pkt.payload[3] = byte(pieceNum)
	pkt.payload[4] = byte(pieceNum >> 8)
	pkt.payload[5] = byte(pieceNum >> 16)
	pkt.payload[6] = byte(pieceNum >> 24)
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

func (tello *Tello) sendFileDone(fID uint16, size int) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	tello.ctrlSeq++
	pkt := newPacket(ptGet, msgFileDone, tello.ctrlSeq, 6)
	pkt.payload[0] = byte(fID)
	pkt.payload[1] = byte(fID >> 8)
	pkt.payload[2] = byte(size)
	pkt.payload[3] = byte(size >> 8)
	pkt.payload[4] = byte(size >> 16)
	pkt.payload[5] = byte(size >> 24)
	tello.ctrlConn.Write(packetToBuffer(pkt))
}
