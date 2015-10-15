package dedup_test

import (
	"testing"

	"bytes"
	"crypto/rand"
	"io"

	"io/ioutil"

	"github.com/klauspost/dedup"
)

func TestWriter(t *testing.T) {
	idx := bytes.Buffer{}
	data := bytes.Buffer{}
	input := &bytes.Buffer{}

	const totalinput = 10 << 20
	_, err := io.CopyN(input, rand.Reader, totalinput)
	if err != nil {
		t.Fatal(err)
	}
	const size = 64 << 10
	b := input.Bytes()
	// Create some duplicates
	for i := 0; i < 50; i++ {
		// Read from 10 first blocks
		src := b[(i%10)*size : (i%10)*size+size]
		// Write into the following ones
		dst := b[(10+i)*size : (i+10)*size+size]
		copy(dst, src)
	}
	input = bytes.NewBuffer(b)
	w, err := dedup.NewFixedWriter(&idx, &data, size)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(w, input)
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}
	removed := ((totalinput) - data.Len()) / size

	t.Log(dedup.BirthDayProblem(totalinput / size))
	t.Log("Index size:", idx.Len())
	t.Log("Data size:", data.Len())
	t.Log("Removed", removed, "blocks")
	// We should get at least 50 blocks
	if removed < 50 {
		t.Fatal("didn't remove at least 50 blocks")
	}
	if removed > 60 {
		t.Fatal("removed unreasonable high amount of blocks")
	}
}

func BenchmarkWriter64K(t *testing.B) {
	input := &bytes.Buffer{}

	const totalinput = 10 << 20
	_, err := io.CopyN(input, rand.Reader, totalinput)
	if err != nil {
		t.Fatal(err)
	}
	const size = 64 << 10
	b := input.Bytes()
	// Create some duplicates
	for i := 0; i < 50; i++ {
		// Read from 10 first blocks
		src := b[(i%10)*size : (i%10)*size+size]
		// Write into the following ones
		dst := b[(10+i)*size : (i+10)*size+size]
		copy(dst, src)
	}
	t.ResetTimer()
	t.SetBytes(totalinput)
	for i := 0; i < t.N; i++ {
		input = bytes.NewBuffer(b)
		w, _ := dedup.NewFixedWriter(ioutil.Discard, ioutil.Discard, size)
		io.Copy(w, input)
		err = w.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkWriter4K(t *testing.B) {
	input := &bytes.Buffer{}

	const totalinput = 10 << 20
	_, err := io.CopyN(input, rand.Reader, totalinput)
	if err != nil {
		t.Fatal(err)
	}
	const size = 8 << 10
	b := input.Bytes()
	// Create some duplicates
	for i := 0; i < 500; i++ {
		// Read from 10 first blocks
		src := b[(i%10)*size : (i%10)*size+size]
		// Write into the following ones
		dst := b[(10+i)*size : (i+10)*size+size]
		copy(dst, src)
	}
	t.ResetTimer()
	t.SetBytes(totalinput)
	for i := 0; i < t.N; i++ {
		input = bytes.NewBuffer(b)
		w, _ := dedup.NewFixedWriter(ioutil.Discard, ioutil.Discard, size)
		io.Copy(w, input)
		err = w.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBirthday(t *testing.T) {
	t.Log("Hash size is", dedup.HashSize*8, "bits")
	t.Log("1GiB, 1KiB blocks:")
	t.Log(dedup.BirthDayProblem((1 << 30) / (1 << 10)))
	t.Log("It will use", dedup.FixedMemUse((1<<30)/(1<<10))>>20, "MiB memory")
	t.Log("1TiB, 4KiB blocks:")
	t.Log(dedup.BirthDayProblem((1 << 40) / (4 << 10)))
	t.Log("It will use", dedup.FixedMemUse((1<<40)/(4<<10))>>20, "MiB memory")
	t.Log("1PiB, 4KiB blocks:")
	t.Log(dedup.BirthDayProblem((1 << 50) / (4 << 10)))
	t.Log("It will use", dedup.FixedMemUse((1<<50)/(4<<10))>>20, "MiB memory")
	t.Log("1EiB, 64KiB blocks:")
	t.Log(dedup.BirthDayProblem((1 << 60) / (64 << 10)))
	t.Log("It will use", dedup.FixedMemUse((1<<60)/(64<<10))>>20, "MiB memory")
	t.Log("1EiB, 1KiB blocks:")
	t.Log(dedup.BirthDayProblem((1 << 60) / (1 << 10)))
	t.Log("It will use", dedup.FixedMemUse((1<<60)/(1<<10))>>20, "MiB memory")
}