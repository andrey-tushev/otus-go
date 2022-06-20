package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()

	data, _ := r.File[0].Open()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(data, "biz")
	}
}
