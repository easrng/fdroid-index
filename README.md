# fdroid-index
Compressed indexes for searching the main F-Droid repo, updated daily at midnight UTC.

## Endpoints
 - `https://easrng.github.io/fdroid-index/index.zst`  
   A zstd-compressed [stork](https://stork-search.net/) index
 - `https://easrng.github.io/fdroid-index/hash.txt`  
   The sha256 hash of index.zst, for update checks. 

## Limitations
 - Indexes are only generated for the main F-Droid repo and not any others (ex. Guardian Project) #1
 - Indexes are only generated for English #2

It wouldn't be very complicated to fix those, I just haven't yet.
