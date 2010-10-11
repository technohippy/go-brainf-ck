/*
BrainCrash (http://cfs.maxn.jp/neta/BrainCrash.php)

How to build:
  $ 6g bc.go
  $ 6l -o braincrash bc.6

How to exec:
  Console:
    $ braincrash
  Read a file:
    $ braincrash [filename]
  Console (brainf*ck mode):
    $ braincrash -f
  Read a file (brainf*ck mode):
    $ braincrash [filename]

How to use:
  Commands (http://wikipedia.org/wiki/Brainfuck):
    Character Meaning
    >         increment the data pointer (to point to the next cell to the right).
    <         decrement the data pointer (to point to the next cell to the left).
    +         increment (increase by one) the byte at the data pointer.
    -         decrement (decrease by one) the byte at the data pointer.
    .         output the value of the byte at the data pointer.
    ,         accept one byte of input, storing its value in the byte at the data pointer.
    [         if the byte at the data pointer is zero, then instead of moving the instruction pointer forward to the next command, jump it forward to the command after the matching ] command*.
    ]         if the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward to the next command, jump it back to the command after the matching [ command*.
    |         or
    &         and
    ~         not
    ^         xor
*/

package main

import (
  "flag"
  "fmt"
  "bufio"
  "os"
)

const (
  PROMPT = "\n%% "
)

var is_brainf_ck = false

func read_line(file *os.File) string {
  line, _ := bufio.NewReader(file).ReadString('\n')
  return line
}

func braincrash(str string) {
  program := []byte(str)
  var memory [30000]byte
  p := 0
  var curr byte

  if !is_brainf_ck {
    ini := "Hello, world!"
    for i := 0; i < len(ini); i++ {
      memory[i] = ini[i]
    }
  }

  for pc := 0; pc < len(program); pc++ {
    switch program[pc] {
    case '>': p++
    case '<': p--
    case '+': memory[p]++
    case '-': memory[p]--
    case '.': fmt.Printf("%c", memory[p])
    case ',': memory[p] = read_line(os.Stdin)[0]
    case '[': 
      if memory[p] == 0 { 
        depth := 0
        for {
          pc++
          if program[pc] == '[' {depth += 1}
          if program[pc] == ']' { if depth == 0 { break } else { depth -= 1} }
        }
      }
    case ']': 
      if memory[p] != 0 { 
        depth := 0
        for {
          pc--
          if program[pc] == ']' {depth += 1}
          if program[pc] == '[' { if depth == 0 { break } else { depth -= 1} }
        }
      }
    }

    if !is_brainf_ck {
      switch program[pc] {
      case '|': curr = memory[p]; p++; memory[p] |= curr
      case '&': curr = memory[p]; p++; memory[p] &= curr
      case '~': memory[p] = ^memory[p]
      case '^': curr = memory[p]; p++; memory[p] ^= curr
      }
    }
  }

  if !is_brainf_ck {
    for i := p; memory[i] != 0; i++ {
      fmt.Printf("%c", memory[i])
    }
  }
}

func read_and_braincrash(file *os.File) bool {
  program := read_line(file)
  if program == "exit\n" { return false }
  braincrash(program)
  return true
}

func main() {
  flag.BoolVar(&is_brainf_ck, "f", false, "brainf*ck mode")
  flag.Parse()
  if flag.NArg() == 0 {
    if is_brainf_ck { fmt.Printf("[BRAINF*CK MODE]\n") }
    fmt.Printf("Type 'exit' to exit")
    for {
      fmt.Printf(PROMPT)
      if !read_and_braincrash(os.Stdin) { break }
    }
  } else {
    filename := flag.Arg(0)
    file, err := os.Open(filename, os.O_RDONLY, 0666)
    if (err == nil) {
      read_and_braincrash(file)
    }
  }
}
