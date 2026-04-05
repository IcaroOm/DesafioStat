package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func main() {
	f, err := os.Open("altura_mundo.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buffer := make([]byte, 4<<20) // 4MB mais que isso fica mais lento

	var (
		count     int64
		mean      float64
		m2        float64
		bytesRead int64
	)


	var num float64
	var frac float64
	var fracDiv float64 = 1
	var parsingFrac bool

	for {
		n, err := f.Read(buffer)
		if n > 0 {
			bytesRead += int64(n)
			for i := 0; i < n; i++ {
				c := buffer[i]
				
				if c == '\r' {  // corrige o bug que aumentou a media
				    continue
				}

				if c == '.' {
					parsingFrac = true
					continue
				}
				if c == '\n' {
					x := num + frac/fracDiv
					count++
					delta := x - mean
					mean += delta / float64(count)
					m2 += delta * (x - mean)
					// reseta o parser
					num = 0
					frac = 0
					fracDiv = 1
					parsingFrac = false
					continue
				}
				digit := float64(c - '0')
				if parsingFrac {
					frac = frac*10 + digit
					fracDiv *= 10
				} else {
					num = num*10 + digit
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	variance := m2 / float64(count)
	stddev := math.Sqrt(variance)

	fmt.Printf("Linhas: %d\n", count)
	fmt.Printf("Media: %.6f\n", mean)
	fmt.Printf("Desvio Padrão: %.6f\n", stddev)
}