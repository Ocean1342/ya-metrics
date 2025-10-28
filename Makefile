pprof_cpu_server:
	go tool pprof -http=":9090" -seconds=30 http://localhost:8080/debug/pprof/profile && \
	curl -sK -v http://localhost:8080/debug/pprof/profile > profile-under-pressure.out && \
    go tool pprof -http=":9090" -seconds=30 ./profiles/server/profile-under-pressure.out

pprof_analyse:
	go tool pprof -http=":9090" -seconds=30 heap.out


pprof_mem_server:
	go tool pprof -http=":9090" -seconds=30 http://localhost:8080/debug/pprof/heap && \
	curl -sK -v http://localhost:8080/debug/pprof/heap > heap-under-pressure_v4.out && \
    go tool pprof -http=":9090" -seconds=30 ./profiles/server/after/heap-under-pressure.out

pprof_diff:
	pprof -top -diff_base=profiles/server/before/heap-under-pressure.out profiles/server/after/heap-under-pressure.out
