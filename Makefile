main.out: *.go
	go build -o aml.out

profile:
	go tool pprof cpu.prof
