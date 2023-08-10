http://localhost:6060/debug/pprof/

go tool pprof http://localhost:6060/debug/pprof/heap

go tool pprof -http=:8080 C:\Users\Developer\pprof\pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.001.pb.gz