package secio

import (
	"sync"
	"testing"

	ci "github.com/libp2p/go-libp2p-crypto"
)

func BenchmarkSessionData1B(b *testing.B)      { benchmarkSessionData(b, 1) }
func BenchmarkSessionData10B(b *testing.B)     { benchmarkSessionData(b, 10) }
func BenchmarkSessionData100B(b *testing.B)    { benchmarkSessionData(b, 100) }
func BenchmarkSessionData1000B(b *testing.B)   { benchmarkSessionData(b, 1000) }
func BenchmarkSessionData10000B(b *testing.B)  { benchmarkSessionData(b, 10000) }
func BenchmarkSessionData100000B(b *testing.B) { benchmarkSessionData(b, 100000) }

func writeData(b *testing.B, sess Session, data []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	rwc := sess.ReadWriter()

	for n := 0; n < b.N; n++ {
		err := rwc.WriteMsg(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func readData(b *testing.B, sess Session, data []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	rwc := sess.ReadWriter()
	dataLen := len(data)

	for n := 0; n < b.N; n++ {
		readData, err := rwc.ReadMsg()
		if err != nil {
			b.Fatal(err)
		}

		if len(readData) != dataLen {
			b.Fatal("Read data length didn't match written.")
		}
	}
}

// Benchmark sending and recieving data on a pair of sessions.
// For this test, we assume setting up the session is trivial.
func benchmarkSessionData(b *testing.B, numBytes int) {
	client_sg := NewTestSessionGenerator(ci.RSA, 1024, b)
	server_sg := NewTestSessionGenerator(ci.RSA, 1024, b)

	client_sess, server_sess := NewTestSessionPair(client_sg, server_sg, b)
	data := make([]byte, numBytes)
	var wg sync.WaitGroup
	wg.Add(2)

	b.ResetTimer()

	go writeData(b, client_sess, data, &wg)
	go readData(b, server_sess, data, &wg)
	wg.Wait()
}

func BenchmarkSessionSetupNoDelay(b *testing.B)           { benchmarkSessionSetup(b, true, true) }
func BenchmarkSessionSetupClientDelay(b *testing.B)       { benchmarkSessionSetup(b, false, true) }
func BenchmarkSessionSetupServerDelay(b *testing.B)       { benchmarkSessionSetup(b, true, false) }
func BenchmarkSessionSetupClientServerDelay(b *testing.B) { benchmarkSessionSetup(b, false, false) }

// Benchmark setting up the session optionally setting the TCP NO_DELAY flags
func benchmarkSessionSetup(b *testing.B, client, server bool) {
	client_sg := NewTestSessionGenerator(ci.RSA, 1024, b)
	server_sg := NewTestSessionGenerator(ci.RSA, 1024, b)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = NewTestSessionPairNoDelay(client_sg, server_sg, client, server, b)
	}
}
