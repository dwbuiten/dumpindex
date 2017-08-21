# ffmsindex/dumpindex

A small package/tool for parsing FFMS2 Version 2 indexes for informational purposes.

# Documentation

See: http://godoc.org/github.com/dwbuiten/dumpindex/ffmsindex

# ffindex Versions

Each supported ffindex version is tagged in the git repository as `vN`, where `N` is the
ffindex version. For example: `v1` or `v4`.

# Tool

This dumpindex tool can be used to dump a given ffindex file to human-readable JSON, via:

```
$ dumpindex file.ffindex
{
    "Header": {
        "ID": 1402079347,
        "Version": {
            "Major": 2,
            "Minor": 22,
            "Micro": 0,
            "Bump": 0
        },
        "IndexVersion": 1,
        "Tracks": 1,
        "Decoder": 1,
        "ErrorHandling": 0,
        "AVUtilVersion": {
            "Major": 55,
            "Minor": 11,
            "Micro": 100
        },

...
```
