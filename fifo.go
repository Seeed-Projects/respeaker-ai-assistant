package main

type FifoBuffer struct {
	buffer []byte
	size   int
	start  int
}

func NewFifoBuffer(capacity int) *FifoBuffer {
	return &FifoBuffer{
		buffer: make([]byte, capacity),
		size:   0,
		start:  0,
	}
}

func (f *FifoBuffer) Write(p []byte) (n int, err error) {
	if len(p) >= len(f.buffer) {
		copy(f.buffer, p[len(p)-len(f.buffer):])
		f.start = 0
		f.size = len(f.buffer)
		return len(p), nil
	}

	writePos := (f.start + f.size) % len(f.buffer)
	if writePos+len(p) <= len(f.buffer) {
		copy(f.buffer[writePos:], p)
	} else {
		part1 := len(f.buffer) - writePos
		copy(f.buffer[writePos:], p[:part1])
		copy(f.buffer[:], p[part1:])
	}

	f.size += len(p)
	if f.size > len(f.buffer) {
		f.start = (f.start + (f.size - len(f.buffer))) % len(f.buffer)
		f.size = len(f.buffer)
	}

	return len(p), nil
}

func (f *FifoBuffer) Read() []byte {
	if f.size == 0 {
		return nil
	}

	data := make([]byte, f.size)
	if f.start+f.size <= len(f.buffer) {
		copy(data, f.buffer[f.start:f.start+f.size])
	} else {
		part1 := len(f.buffer) - f.start
		copy(data[:part1], f.buffer[f.start:])
		copy(data[part1:], f.buffer[:f.size-part1])
	}

	return data
}
