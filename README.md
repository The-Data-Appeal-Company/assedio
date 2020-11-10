# Assedio

[![Go Report Card](https://goreportcard.com/badge/github.com/The-Data-Appeal-Company/assedio)](https://goreportcard.com/report/github.com/The-Data-Appeal-Company/assedio)
![Go](https://github.com/The-Data-Appeal-Company/assedio/workflows/Go/badge.svg)
[![license](https://img.shields.io/github/license/The-Data-Appeal-Company/assedio.svg)](LICENSE)

![alt Assedio di Firenze, Giovanni Stradano, Firenze, 1530](https://upload.wikimedia.org/wikipedia/commons/2/2b/Siege_of_Florence.JPG)

Assedio di Firenze, Giovanni Stradano, Firenze, 1530

### Simple concurrent http calls tool
Assedio is a tool to make concurrency http calls read from file.
Despite the commons http calls tool such as `siege`,
 Assedio is intended to call each http endpoint once, so that you can reproduce the same test for performance comparision
 
### Usage

```
assedio fight -k 10 -f test_file
```