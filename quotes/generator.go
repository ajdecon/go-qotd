package quotes

import(
    "io/ioutil"
    "log"
    "math/rand"
    "strings"
)

// Maximum allowed length for QOTD server is 512 as defined in
// RFC 865: http://tools.ietf.org/html/rfc865
const MaxLength = 512

// Return a channel which feeds quotes based on a file
func FileGenerator(filename string) chan string {
    qchannel := make(chan string, 5)
    go generateQuotes(filename, qchannel)
    return qchannel
}

// Build a cache of quotes in a map so we can select randomly by
// line number.

func cacheQuotes(filename string) ([]string, int) {
    contents, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Could not read in %s\n", filename)
    }
    quotes := strings.Split(string(contents), "\n%\n")

    // Trim the quotes
    for i := range quotes {
        if len(quotes[i]) > MaxLength {
            quotes[i] = string([]byte(quotes[i])[0:MaxLength-3]) // MaxLength - 3 for ellipses
            quotes[i] = quotes[i] + "..."
        }
    }

    return quotes, len(quotes)
}

// Generate quotes
func generateQuotes(filename string, qchannel chan string) {
    cache, count := cacheQuotes(filename)
    for {
        qchannel <- cache[rand.Intn(count)]
    }
}
