package SAXPY_Bench

import (
    "testing"
)

func TestXYZ(t *testing.T) {

}

func BenchmarkSaxpySerial(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Saxpy_serial()
    }
}

func BenchmarkSaxpyParallel(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Saxpy_parallel()
    }
}

func BenchmarkSaxpyImprove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Saxpy_improve()
    }
}
