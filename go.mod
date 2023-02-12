module fdroid-index

go 1.19

replace github.com/avast/apkverifier => github.com/easrng/fdroid-repo-verifier v0.0.0-20230212071853-3dbdcc239ba3

require github.com/avast/apkverifier v0.0.0-20221110131049-7720fc1ebef0

require (
	github.com/avast/apkparser v0.0.0-20221110131626-bc2b7ccc9d3e // indirect
	github.com/itchyny/gojq v0.12.11 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/klauspost/compress v1.15.12 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
)
