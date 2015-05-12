package Dot_Bench

import (
    "testing"
)

func TestXYZ(t *testing.T) {

}

func BenchmarkDotSerialA(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Dot_serial_A()
    }
}

func BenchmarkDotParallelA(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Dot_parallel_A()
    }
}

func BenchmarkDotImproveA(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Dot_improve_A()
    }
}

func BenchmarkDotSerialB(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Dot_serial_B()
    }
}

func BenchmarkDotParallelB(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Dot_parallel_B()
    }
}

func BenchmarkDotImproveB(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Dot_improve_B()
    }
}
