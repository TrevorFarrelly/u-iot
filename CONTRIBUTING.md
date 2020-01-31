# Contributing to u-iot

Contributions are always welcome. Check out the project [milestones](https://github.com/TrevorFarrelly/u-iot/projects) or open [issues](https://github.com/TrevorFarrelly/u-iot/issues) for good places to start.

### Code of Conduct

Collaboration around u-iot must be done in accordance to our [Code of Conduct](CODE-OF-CONDUCT.md).

### Getting Started

Great! You want to contribute. Please follow this basic outline when making changes:
1. Fork the repository.
2. Make the changes in a feature branch on your fork.
3. Run [gofmt](#Coding-Style).
4. Open a PR!

### PR/Issue Tags

When opening an issue or a pull request, please tag it with a relevant label. Examples include *doc*, *bug*, *feature*, etc. If the pull request is fixing a bug or implementing a feature that has an issue associated with it, add
>fixes #(*issue number*)

to the body of the PR. The issue will automatically close when the PR is merged.

### New Features

u-iot at its core is meant to be a lightweight library that can easily be included on small SD card or embedded device. That said, feel free to pitch the idea using the *feature request* issue tag before opening a PR.

### Developer Sign-Off

Please use git's `-s` flag or append each nontrivial git commit with the following line when contributing.
```
Signed-off-by: Your Name <username@youremail.com>
```
A nontrivial git commit is one with enough code changes that it could be considered intellectual property. Please only sign off on the commit if you can attest to the [Developer's Certificate of Origin](https://developercertificate.org/).

### Coding Style

As with most Go programs, u-iot follows the [Effective Go](https://golang.org/doc/effective_go.html) style guide. Therefore, `gofmt` and `golint` know what is best. Write your code however you please, then run it through `gofmt` before submitting.
