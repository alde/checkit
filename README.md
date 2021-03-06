# Checkit
[![Go Report Card](https://goreportcard.com/badge/github.com/alde/checkit)](https://goreportcard.com/report/github.com/alde/checkit)

Combine checkstyle and spotbugs (and findbugs) reports into one json file.

## Why?

If you're hosting CI infrastructure for different teams but want to consume reports downstream, restricting teams to use tools that can report in a single format can be hard or limiting.
If you can convert these reports into a common format, you can analyze build reports in an easier fashion.

## Usage

```
checkit -exclude testdata,fixtures -output combined-checkit.json
```

## Installation

```
go get github.com/alde/checkit
```
