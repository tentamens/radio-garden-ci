package helpers

import (
	"io"
	"time"
)

type ProcessorReader struct {
	Reader io.Reader
	Buffer []byte
	Func   func([]byte) // Your function to process the PCM data chunk
}

// NewProcessorReader initializes the reader with the correct buffer size
func NewProcessorReader(r io.Reader, duration time.Duration, sampleRate, channels, bytesPerSample int, processFunc func([]byte)) *ProcessorReader {
	bytesPerSecond := sampleRate * channels * bytesPerSample
	bufferSize := int(float64(bytesPerSecond) * duration.Seconds())

	// Ensure the buffer size is even if you're working with 16-bit data
	if bufferSize%2 != 0 {
		bufferSize++
	}

	return &ProcessorReader{
		Reader: r,
		Buffer: make([]byte, bufferSize),
		Func:   processFunc,
	}
}

// Read implements the io.Reader interface.
func (pr *ProcessorReader) Read(p []byte) (n int, err error) {
	// 1. Read data from the underlying FFmpeg pipe into the buffer
	bytesRead, err := io.ReadFull(pr.Reader, pr.Buffer)
	if err != nil && err != io.ErrUnexpectedEOF {
		// If an EOF happens mid-read, bytesRead will contain partial data,
		// which might still be useful. io.ReadFull returns io.ErrUnexpectedEOF
		// if it can't fill the buffer completely.
		if err == io.ErrUnexpectedEOF && bytesRead > 0 {
			// Process the partial buffer then return the error
			pr.Func(pr.Buffer[:bytesRead])
			// Fallthrough to step 2/3
		} else {
			return 0, err
		}
	}

	// 2. Call the processing function on the filled buffer
	pr.Func(pr.Buffer[:bytesRead])

	// 3. Copy the buffer content to the destination slice 'p' (for oto)
	// The oto library will call Read multiple times to fill its internal buffer.
	// To ensure oto is supplied with *at least* the size of the ProcessorReader's
	// buffer (8820 bytes), we need to fill 'p' which oto supplies (often larger).
	// We only ever successfully copy the size of our buffer (bytesRead) or the
	// size of 'p', whichever is smaller.

	copySize := bytesRead
	if len(p) < copySize {
		copySize = len(p)
	}
	copy(p[:copySize], pr.Buffer[:copySize])

	// Return the number of bytes successfully copied to 'p'
	return copySize, nil
}
