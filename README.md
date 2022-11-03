# Self Organising Multi Agent Systems 2022

## [Specifications](spec.pdf)

## Project Structure

Following some Golang Standards [[1]](https://github.com/golang-standards/project-layout), [[2]](https://medium.com/sellerapp/golang-project-structuring-ben-johnson-way-2a11035f94bc) and previous year project structures ([2021](https://github.com/SOMAS2021/SOMAS2021))

```
.
├── cmd
│   └── (Executable Outputs)
├── config.json
├── docs
│   └── (Documentation Files)
├── frontend
│   └── (Frontend Implementation)
├── pkg
│   ├── agents
│   │   └── (Teams Agents)
│   └── infra
│       └── (Infrastructure Implementation)
└── scripts
    └── (Automation/Execution scripts)
```

## Programming Language

* Download [GoLang](https://go.dev/dl/) 1.16 (?)

### Get Started

* [GoLang Get Started](https://go.dev/learn/)

## Project Management

* Use [Github Issues](https://github.com/features/issues) to track tasks, deadlines, assignees, project progress etc.

## Contribution Guidelines

A lot of these guidelines are from the [SOMAS2020 repo](https://github.com/SOMAS2020/SOMAS2020/blob/main/docs/SETUP.md) :)

### Coding Rules

1. You're encouraged to use [VSCode](https://code.visualstudio.com/) with the [Go extension](https://code.visualstudio.com/docs/languages/go).
2. Trust the language server. Red lines == death. Yellow lines == close to death. An example where it might be very tempting to let yellow lines pass are in `struct`s:
```golang
type S struct {
    name string
    age  int
}
s1 := S {"pitt", 42} // NO
s2 := S {"pitt"} // NO (even though it initialises age to 0)
s3 := S{name: "pitt", age: 42} // OK, if we add fields into S, this will still be correct
s4 := S{name: "pittson"} // OK if `pittson`'s age is 0
```
3. Write tests where required. There are many guides online on how to do this in Golang. Tests will be run alongside CI when you pull into the main repo. If anyone breaks anything, it's easy to observe that if you have tests. Otherwise, your code will be broken unknowingly.

4. DO NOT TOUCH code you don't own unless you have a good reason to. If you have a good reason to, do it in a separate PR and notify the owners of the code.

5. Do not use `panic` or `die` - return an `error` instead!

6. Do not use system-specific packages (e.g. `internal/syscall/unix`).

7. Use the superior `errors.Errorf` to create your errors so that we have a stack trace.
### Code Reviews and PRs
- Do not push to the `main` branch. Make your own branch and create a PR into `main` when it's ready for review.
- When working on your own team's features, please name your branch as: `teamX-FEATURE_NAME-WHATEVER_YOU_LIKE_HERE`
- Do not use force push. Use `git push --force-with-lease` instead.
- When ready to merge into your team's feature branch, create a PR to merge into `teamX-FEATURE_NAME`. When the feature is complete, then create a PR into `main`.
- Make sure that you have reviewed your own code before creating the PR.
- Keep PRs small: if they are too large, you will be told to split your code into smaller PRs. This ensures they can be reviewed properly.
- You need to make sure your code is up-to-date with the `main` branch. Merge commits are *not allowed*: learn how to rebase. Ask someone in the #git-env-infra channel if you don't know how to do this.
- Do not review your own code. That completely defeats the purpose of code review.
- Team leads: when doing a PR for your team's code, do not approve it yourself - get another team lead to review it.
- Review PRs in a timely manner! Ideally by the next day so that other teams aren't blocked.
- If you create a PR, use the "assign" feature on the PR to assign who should be merging it once the review is completed (this can be you)
- If you are a reviewer, do not merge in PRs that are not assigned to you once you finish reviewing.
