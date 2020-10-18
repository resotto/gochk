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

- ZERO Dependency (Only Standard Library)
- Simple and Easy to Read Results

---

## Table of Contents

- [Getting Started](#getting-started)
- [How to use gochk](#how-to-use-gochk)
- [How to see results](#how-to-see-results)
- [Performance](#performance)
- [Customization](#customization)
- [](#)
- [](#)

## Getting Started

```zsh
go get -u github.com/resotto/gochk
```

Please edit `dependencyOrders` in `gochk/configs/config.json` according to your dependency rule, whose smaller index value means outer circle.

```json
"dependencyOrders": ["external", "adapter", "application", "domain"],
```

And then, let's gochk it!

```zsh
go run cmd/gochk/main.go {CheckTargetPath}
```

## How to use gochk

## How to see results

## Performance

## Customization

## Build

How to build ...

## Feedback

- [Feel free to write your thoughts](https://github.com/resotto/gochk/issues/1)
- Report a bug to [Bug report](https://github.com/resotto/gochk/issues/2).

## License

[GNU General Public License v3.0](https://github.com/resotto/gochk/blob/master/LICENSE).
