Team5 folder structure:


```
.
Team5
├── agent5.go
│   └── (Executable signature functions)
├── commons5
│   └── (common structures&functions called for various parts)
├── allocation
│   └── (lootallocation functions Wenlin&Matthew)
├── leaderFight
│   └── (fight proposal functions)
.
.
.
<<<<<<< HEAD
├── cmd
│   └── (Executable Outputs)
├── docs
│   └── (Documentation Files)
├── web
│   └── (Frontend/Backend Implementation)
├── pkg
│   └── infra
│       └── (Infrastructure Implementation)
|       └── teams
|           └── (Individual Team Agents/Experiments)
├── .env (Environmental variables for Infrastructure)
└── scripts
    └── (Automation/Execution scripts)
=======
Team5
├── agent5.go
│   └── (Executable signature functions)
├── commons5
│   └── (common structures&functions called for various parts)
├── allocation
│   └── (lootallocation functions Wenlin&Matthew)
├── leaderFight
│   └── (fight proposal functions)
.
.
.

## Project Management

* Use [Github Issues](https://github.com/features/issues) to track tasks, deadlines, assignees, project progress etc.

## Contribution Guidelines

A lot of these guidelines are from the [SOMAS2021](https://github.com/SOMAS2021/SOMAS2021/blob/main/README.md)
and [SOMAS2020](https://github.com/SOMAS2020/SOMAS2020/blob/main/docs/SETUP.md) repos :)

### Coding Rules

1. You're encouraged to use [VSCode](https://code.visualstudio.com/) with
   the [Go extension](https://code.visualstudio.com/docs/languages/go).
2. Trust the language server. Red lines == death. Yellow lines == close to death. An example where it might be very
   tempting to let yellow lines pass are in `struct`s:

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
