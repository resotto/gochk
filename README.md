<h1 align="center">go chk</h1>

<p align="center">
  <a href="https://github.com/resotto/gochk/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-GPL%20v3.0-brightgreen.svg" /></a>
</p>

<p align="center">
  Static Dependency Rule Analysis Tool for Go Files
</p>

<p align="center">
  <img src="https://user-images.githubusercontent.com/19743841/96338043-67719c80-10c6-11eb-9a5f-3b672356a9d6.gif">
</p>

---

What is gochk?

- gochk checks for .go files' [Clean Architecture The Dependency Rule](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html#the-dependency-rule) and prints its results.

  > This rule says that source code dependencies can only point inwards. Nothing in an inner circle can know anything at all about something in an outer circle.

Why gochk?

- ZERO Dependency
- Simple and Easy to Read Outputs

---

## Table of Contents

- [Getting Started](#getting-started)
- [How to use gochk](#how-to-use-gochk)
- [Configuration](#configuration)
- [How to see results](#how-to-see-results)
- [Performance](#performance)
- [Testing](#testing)
- [Customization](#customization)
- [Build](#build)
- [Feedback](#feedback)
- [License](#license)

## Getting Started

```zsh
go get -u github.com/resotto/gochk
cd ${GOPATH}/src/github.com/resotto/gochk
```

Please edit paths defined by `dependencyOrders` in `gochk/configs/config.json` according to your dependency rule, whose smaller index value means outer circle.

```json
"dependencyOrders": ["external", "adapter", "application", "domain"],
```

And then, let's gochk your target path!

```zsh
go run cmd/gochk/main.go {CheckTargetPath}
```

If you have [Goilerplate](https://github.com/resotto/goilerplate), you can also gochk it like:

```zsh
go run cmd/gochk/main.go ../goilerplate
```

## How to use gochk

## How to see results

### Quick Result

You can check whether there are violations or not quickly by checking the end of results.

If you see the following text, there are violations!

```
2020/10/19 23:37:03 Dependencies which violate dependency orders found!
exit status 1
```

If you can see the following AA at the bottom, there are no violations, Congrats🎉

```
2020/10/19 23:57:25 No violations
    ________     _______       ______    __     __    __   _ _
   /  ______\   /  ___  \     /  ____\  |  |   |  |  |  | /   /
  /  /  ____   /  /   \  \   /  /       |  |___|  |  |  |/   /
 /  /  |_   | |  |     |  | |  |        |   ___   |  |      /
 \  \    \  | |  |     |  | |  |        |  |   |  |  |  |\  \
  \  \___/  /  \  \___/  /   \  \_____  |  |   |  |  |  | \  \
   \_______/    \_______/     \_______\ |__|   |__|  |__|  \__\
```

### Result types

gochk displays results with colors by default:

- None: ![#008080](https://via.placeholder.com/15/008080/000000?text=+) <span style="color:teal">teal</span>
  - which means no dependencies (no imports).
- Verified: ![#008000](https://via.placeholder.com/15/008000/000000?text=+) <span style="color:green">grean</span>
  - which means there are dependencies with no violation.
- Ignored: ![#FFFF00](https://via.placeholder.com/15/FFFF00/000000?text=+) <span style="color:yellow">yellow</span>
  - which means the path is ignored (not checked) by gochk .
- Warning: ![#800080](https://via.placeholder.com/15/800080/000000?text=+) <span style="color:purple">purple</span>
  - which means something happened when gochk checked it.
- Violated: ![#FF0000](https://via.placeholder.com/15/FF0000/000000?text=+) <span style="color:red">red</span>
  - which means there are dependencies which violates dependency rule.

For `none`, `verified`, and `ignored`, they display only the file path.

```
[None]     ../goilerplate/internal/app/adapter/postgresql/conn.go
```

```
[Verified] ../goilerplate/cmd/app/main.go
```

```
[Ignored]  ../goilerplate/.git
```

For `warning`, it displays what happened to the file.

```
[Warning]  open /Users/resotto/go/src/github.com/resotto/goilerplate/internal/app/application/usecase/lock.go: permission denied
```

For `violated`, it displays the file path, its dependency, and how it violates.

```
[Violated] ../goilerplate/internal/app/domain/temp.go imports "github.com/resotto/goilerplate/internal/app/adapter/postgresql/model"
 => domain depends on adapter
```

## Configuration

`gochk/configs/config.json` has configuration values.

```json
{
  "targetPath": ".",
  "dependencyOrders": ["external", "adapter", "application", "domain"],
  "ignore": ["test", "_test", ".git"],
  "printViolationsAtTheBottom": false
}
```

- `targetPath` is default target path of gochk.
  - If you specify this parameter, you don't need to pass command line argument.
- `dependencyOrders` are the paths of each circles in Clean Architecture.

  - For example, if you have following four circles such as Clean Architecture:

    - "External" (most outer)
    - "Adapter"
    - "Application"
    - "Domain" (most inner, the core)

    you should specify them from the outer to the inner like `["external", "adapter", "application", "domain"]`.

  - If you have other layered architecture, you could specify its layers to this parameter as well.

- `ignore` is the path ignored by gochk which can be file name or dir name.
- `printViolationsAtTheBottom` is how gochk prints results about violations of the dependency rule.

  - If `true`, you can see violations at the bottom like:
  <p align="center">
    <img src="https://user-images.githubusercontent.com/19743841/96462260-58beed00-1260-11eb-8459-4938d184cb37.gif">
  </p>

  - If `false`, you see them disorderly:
  <p align="center">
    <img src="https://user-images.githubusercontent.com/19743841/96462316-670d0900-1260-11eb-9b3c-882a3f7adf93.gif">
  </p>

## Performance

## Testing

### Test Package Structure

```zsh
├── internal
│   └── gochk
│       ├── calc_internal_test.go # Unit test (internal)
│       └── read_internal_test.go # Unit test (internal)
└── test
    ├── data                      # Test data
    └── performance
        └── check_test.go         # Performance test
```

### Unit Test

```zsh
cd internal
```

```zsh
~/go/src/github.com/resotto/gochk/internal/gochk (master) > go test ./... # Please specify -v if you need detailed outputs
ok      github.com/resotto/gochk/internal/gochk (cached)
```

You can also clean test cache with `go clean -testcache`.

```zsh
~/go/src/github.com/resotto/gochk/internal/gochk (master) > go clean -testcache
~/go/src/github.com/resotto/gochk/internal/gochk (master) > go test ./...
ok      github.com/resotto/gochk/internal/gochk 0.065s # Not cache
```

### Performance Test

Performance test will take few minutes.

```zsh
cd test/performance
```

```zsh
~/go/src/github.com/resotto/gochk/test/performance (master) > go test ./...
ok      github.com/resotto/gochk/test/performance       64.661s
```

## Customization

## Build

How to build ...

## Feedback

- [Feel free to write your thoughts](https://github.com/resotto/gochk/issues/1)
- Report a bug to [Bug report](https://github.com/resotto/gochk/issues/2).

## License

[GNU General Public License v3.0](https://github.com/resotto/gochk/blob/master/LICENSE).
