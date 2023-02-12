# fdroid-index
Compressed indexes for searching the main F-Droid repo, updated daily at midnight UTC.

## Endpoints
 - `https://easrng.github.io/fdroid-index/index.zst`  
   A zstd-compressed [stork](https://stork-search.net/) index
 - `https://easrng.github.io/fdroid-index/hash.txt`  
   The sha256 hash of index.zst, for update checks. 
