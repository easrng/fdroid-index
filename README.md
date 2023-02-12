# fdroid-index
Compressed indexes for searching the main F-Droid repo, updated daily at midnight UTC.

## Endpoints
 - [`/index.zst`](https://easrng.github.io/fdroid-index/index.zst)  
   A zstd-compressed [stork](https://stork-search.net/) index
 - [`/hash.txt`](https://easrng.github.io/fdroid-index/hash.txt)  
   The sha256 hash of index.zst, for update checks. 

## Limitations
 - https://github.com/easrng/fdroid-index/issues/1 Indexes are only generated for the main F-Droid repo and not any others (ex. Guardian Project)
 - https://github.com/easrng/fdroid-index/issues/2 Indexes are only generated for English

It probably wouldn't be very complicated to fix those, I just haven't yet.
