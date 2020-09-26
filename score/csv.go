package score

import (
	"encoding/csv"
	"fmt"
	"go-snake-ai/input"
	"log"
	"os"
	"path"
	"strconv"
)

func NewCSV(writeDir string) *CSV {
	return &CSV{
		writeDir: writeDir,
	}
}

type CSV struct {
	writeDir string
}

func (w *CSV) Write(score int, maxScore int, input input.Input) {
	fileName := fmt.Sprintf("scores_%s_%d.csv", input.Name(), maxScore)
	filePath := path.Join(w.writeDir, fileName)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("failed to write score: %s\n", err)
		return
	}
	writer := csv.NewWriter(f)
	record := []string{
		strconv.Itoa(score),
	}
	writer.Write(record)
	writer.Flush()
}
