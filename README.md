decryption spike
================

1. Have a big CSV file in `./data/bigdata.csv` (example used 3GB)
2. Create an `output` directory
3. Run `time go run encrypt/main.go`
4. Run `time go run decrypt/main.go`

To test it without concurrency (assuming `bigdata.csv` has less than 50m entries):

1. Run `time BATCH_SIZE=50000000 go run encrypt/main.go`
2. Run `time go run decrypt/main.go`

Example decrypt output:

```bash
# Decrypting 1 x 3GB file
➜ time go run decrypt/main.go
go run decrypt/main.go  123.33s user 1.78s system 101% cpu 2:03.31 total

# Decrypting 21 x 150M files (concurrency: 4)
➜ time go run decrypt/main.go
go run decrypt/main.go  188.23s user 3.67s system 365% cpu 52.492 total

# Decrypting 21 x 150M files (concurrency: 8)
➜ time go run decrypt/main.go
go run decrypt/main.go  263.04s user 3.40s system 591% cpu 45.037 total

# Decrypting 21 x 150M files (concurrency: 12)
➜ time go run decrypt/main.go
go run decrypt/main.go  281.83s user 3.68s system 640% cpu 44.560 total
```

### Licence

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

