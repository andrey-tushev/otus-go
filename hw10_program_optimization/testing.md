# Использование CPU и RAM

## Начальные 
    stats_optimization_test.go:45: time used: 653.244035ms / 300ms
    stats_optimization_test.go:46: memory used: 292Mb / 30Mb

    cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
    BenchmarkGetDomainStat
    BenchmarkGetDomainStat-12    	  877502	      1285 ns/op

## После оптимизации
    stats_optimization_test.go:46: time used: 145.435453ms / 300ms
    stats_optimization_test.go:47: memory used: 4Mb / 30Mb

    cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
    BenchmarkGetDomainStat
    BenchmarkGetDomainStat-12    	  877502	      1285 ns/op

# Профилирование

    go test -v -count=1 -timeout=30s -tags bench . -cpuprofile=cpu.out -memprofile=mem.out
    
    go tool pprof -http=":8090" hw10_program_optimization.test mem.out   
    go tool pprof -http=":8090" hw10_program_optimization.test cpu.out

