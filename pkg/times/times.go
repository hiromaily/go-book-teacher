package times

import (
	"fmt"
	"time"
)

func Track(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
