# Contribution Guide

## Prerequisite

**Please note that this project is released with a [Contributor Code of Conduct](#https://github.com/resotto/gochk/blob/master/CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.**

## Where to report violations.

Please email me resotto3@gmail.com.

## How to contribute Gochk.

I appreciate your help.

1. First, please write your issue(proposal) in [Gochk issues](https://github.com/resotto/gochk/issues).

1. Secondly, you must have the following tools and settings on your IDE:

   - `godoc` for docs
   - `goreturns` for format
   - `golint` for lint
   - build on save for `package`
   - lint on save for `package`
   - vet on save for `package`

1. After satisfing the above, please make a branch with `{ISSUE_NUMBER}.{SUMMARY}`.

1. You MUST also fix/add unit tests of your implementation in `internal/gochk/xxx_internal_test.go`.

1. Finally, please make a pull request of it.

If you contribute documents, step 2 & 4 above might be skipped.
