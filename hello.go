package main

import (
    "example.com/hello/raytracer"
    "fmt"
)

import "rsc.io/quote"

func main() {
    raytracer.GenerateImage()
    fmt.Println(quote.Glass())
}
