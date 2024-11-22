# Rust vs. Go

Comparison of Rust and Go in terms of performance.

## Rust

```bash
     ✗ OK
      ↳  99% — ✓ 1865605 / ✗ 13
     ✗ Similarity Returned
      ↳  99% — ✓ 1865605 / ✗ 13

     checks.........................: 99.99%  3731210 out of 3731236
     data_received..................: 310 MB  645 kB/s
     data_sent......................: 38 GB   78 MB/s
     http_req_blocked...............: avg=5.28µs  min=1.38µs   med=4.22µs  max=15.33ms p(90)=5.68µs   p(95)=6.31µs  
     http_req_connecting............: avg=191ns   min=0s       med=0s      max=12.04ms p(90)=0s       p(95)=0s      
   ✓ http_req_duration..............: avg=90.77ms min=972.39µs med=74.79ms max=1m0s    p(90)=161.72ms p(95)=219.01ms
       { expected_response:true }...: avg=90.35ms min=972.39µs med=74.79ms max=59.9s   p(90)=161.71ms p(95)=219ms   
   ✓ http_req_failed................: 0.00%   13 out of 1865618
     http_req_receiving.............: avg=59.68µs min=0s       med=34.8µs  max=22.07ms p(90)=48.32µs  p(95)=68.42µs 
     http_req_sending...............: avg=98.37µs min=29.7µs   med=76.18µs max=20.57ms p(90)=101.4µs  p(95)=142.3µs 
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=90.61ms min=900.59µs med=74.64ms max=1m0s    p(90)=161.55ms p(95)=218.79ms
     http_reqs......................: 1865618 3886.688951/s
     iteration_duration.............: avg=91.7ms  min=1.33ms   med=75.69ms max=1m0s    p(90)=162.66ms p(95)=219.97ms
     iterations.....................: 1865618 3886.688951/s
     vus............................: 1       min=1                  max=800
     vus_max........................: 800     min=800                max=800


running (8m00.0s), 000/800 VUs, 1865618 complete and 0 interrupted iterations
default ✓ [======================================] 000/800 VUs  8m0s
```

## Go

```bash
     ✓ OK
     ✓ Similarity Returned

     checks.........................: 100.00% 1920966 out of 1920966
     data_received..................: 159 MB  332 kB/s
     data_sent......................: 19 GB   40 MB/s
     http_req_blocked...............: avg=5.7µs    min=1.49µs  med=4.12µs  max=19.21ms p(90)=5.68µs   p(95)=6.33µs  
     http_req_connecting............: avg=478ns    min=0s      med=0s      max=10.48ms p(90)=0s       p(95)=0s      
   ✓ http_req_duration..............: avg=177.35ms min=1.93ms  med=22.28ms max=5.25s   p(90)=651.66ms p(95)=905.77ms
       { expected_response:true }...: avg=177.35ms min=1.93ms  med=22.28ms max=5.25s   p(90)=651.66ms p(95)=905.77ms
   ✓ http_req_failed................: 0.00%   0 out of 960483
     http_req_receiving.............: avg=65.69µs  min=12.43µs med=35.92µs max=31.08ms p(90)=51.72µs  p(95)=69.33µs 
     http_req_sending...............: avg=111.55µs min=32.45µs med=75.79µs max=27.21ms p(90)=102.84µs p(95)=164.76µs
     http_req_tls_handshaking.......: avg=0s       min=0s      med=0s      max=0s      p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=177.17ms min=899.7µs med=22.08ms max=5.25s   p(90)=651.49ms p(95)=905.59ms
     http_reqs......................: 960483  2000.998964/s
     iteration_duration.............: avg=178.31ms min=2.38ms  med=23.32ms max=5.25s   p(90)=652.74ms p(95)=906.96ms
     iterations.....................: 960483  2000.998964/s
     vus............................: 1       min=1                  max=800
     vus_max........................: 800     min=800                max=800


running (8m00.0s), 000/800 VUs, 960483 complete and 0 interrupted iterations
default ✓ [======================================] 000/800 VUs  8m0s
```
