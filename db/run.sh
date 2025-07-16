#!/bin/bash
go test -c -gcflags="all=-N -l" -o perf_analysis/testbin

cd perf_analysis

echo "cache L1 test"
perf stat -e L1-dcache-loads,L1-dcache-load-misses ./testbin 
echo "-------------------------------------------------------"

echo "cache L2 test"
perf stat -e LLC-loads,LLC-load-misses ./testbin
echo "-------------------------------------------------------"

echo "branch prediciton test"
perf stat -e branches,branch-misses ./testbin
echo "-------------------------------------------------------"

echo "instruction and cycles test"
perf stat -e instructions,cycles ./testbin
echo "-------------------------------------------------------"

echo "all caches in one hit"
perf stat -e cache-references,cache-misses,L1-dcache-loads,L1-dcache-load-misses,LLC-loads,LLC-load-misses ./testbin
echo "-------------------------------------------------------"

echo "cpu backend bottlenecks"
perf stat -e cycles,instructions,cache-references,cache-misses,branches,branch-misses ./testbin
echo "-------------------------------------------------------"

echo "perf rec cycles"
perf record -e cycles ./testbin
perf report
echo "-------------------------------------------------------"


