package errcode

import "testing"

func Test_getAppErrOccurredInfo(t *testing.T) {

	s := getAppErrOccurredInfo(1)
	t.Log(s)
}

func BenchmarkLastTwoParts(b *testing.B) {
    s := "p/q/r/s/t/u/v/w/x/y/z"
    for i := 0; i < b.N; i++ {
        lastTwoParts(s)
    }
}

func BenchmarkLastTwoParts2(b *testing.B) {
    s := "p/q/r/s/t/u/v/w/x/y/z"
    for i := 0; i < b.N; i++ {
        lastTwoParts2(s)
    }
}