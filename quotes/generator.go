package quotes

import(
    "bufio"
    "os"
    "log"
    "math/rand"
)

// Return a channel which feeds quotes based on a file
func FileGenerator(filename string) chan string {
    qchannel := make(chan string, 5)
    go generateQuotes(filename, qchannel)
    return qchannel
}

// Build a cache of quotes in a map so we can select randomly by
// line number.
func cacheQuotes(filename string) (cache map[int]string, count int) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatalf("Could not open file: %s\n", filename)
    }

    fileScanner := bufio.NewScanner(file)
    count = 0
    cache = make(map[int]string)

    for fileScanner.Scan() {
        cache[count] = fileScanner.Text()
        count += 1
    }

    file.Close()
    return
}

// Generate quotes
func generateQuotes(filename string, qchannel chan string) {
    cache, count := cacheQuotes(filename)
    for {
        qchannel <- cache[rand.Intn(count)]
    }
}
