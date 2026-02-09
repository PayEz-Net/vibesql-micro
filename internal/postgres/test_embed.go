package postgres

import (
    "fmt"
)

func TestEmbed() {
    entries, err := embeddedPostgres.ReadDir("embed")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Embedded files:")
    for _, e := range entries {
        info, _ := e.Info()
        fmt.Printf("  %s (%d bytes)\n", e.Name(), info.Size())
    }
}
