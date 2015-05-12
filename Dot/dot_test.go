package Dot

import "testing"

func TestXYZ(t *testing.T) {

}

func BenchmarkDotSerial(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Dot_serial()
    }
}

func BenchmarkDotParallel(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Dot_parallel()
    }
}
