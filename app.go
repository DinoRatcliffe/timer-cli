package main

import (
    "strings"
    "flag"
    "bytes"
    "os/exec"
    "os"
    "time"
    "fmt"
    "math"
    "log"

    "github.com/google/goterm/term"
    "github.com/faiface/beep/mp3"
    "github.com/faiface/beep/speaker"
    "golang.org/x/crypto/ssh/terminal"
)

func main() {
    // handling the command line arguments
    timerLengthFlag := flag.Int("legth", 10, "Length of the timer in seconds.")
    flag.Parse()

    var timerLength = *timerLengthFlag

    // get terminal size
    width, height, err := terminal.GetSize(0)
    if err != nil {
        log.Fatal(err)
    }
    go update(time.Now(), float64(timerLength), width, height)
    time.Sleep(time.Duration(timerLength + 5) * time.Second)
    cmd :=exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()

}

func timeString(t time.Duration) string {
    days := t / (24 * time.Hour)
    hours := t % (24 * time.Hour)
    mins := hours % time.Hour
    seconds := math.Mod(mins.Seconds(), 60)

    var buffer bytes.Buffer
    if days > 0 {
        buffer.WriteString(fmt.Sprintf("%dd", days))
    }

    if hours/time.Hour > 0 {
        buffer.WriteString(fmt.Sprintf("%dh", hours/time.Hour))
    }

    if mins/time.Minute > 0 {
        buffer.WriteString(fmt.Sprintf("%dm", mins/time.Minute))
    }

    if seconds > 0 {
        buffer.WriteString(fmt.Sprintf("%.0fs", seconds))
    }
    return buffer.String()
}

func update(t time.Time, timerLength float64, width int, height int) {
    played := false
    for range time.Tick(time.Millisecond * 200) {
        // Clear the screen
        cmd :=exec.Command("clear")
        cmd.Stdout = os.Stdout
        cmd.Run()

        elapsed := time.Since(t)

        // calucate time and progress
        currentWidth := int(float64(width) * (float64(elapsed / time.Second) / timerLength))
        if currentWidth > width {
            currentWidth = width
        }

        // construct status strings
        status := "Elapsed: " + timeString(elapsed)
        timerString := timeString(time.Duration(timerLength) * time.Second) + " Timer"
        remaining := "Remaining: " + timeString((time.Duration(timerLength) * time.Second) - elapsed)

        // construct status bar
        line := strings.Repeat("#", currentWidth)
        page := strings.Repeat(line+"\n", height-2)
        pad := strings.Repeat(" ", ((width - len(status) - len(timerString) - len(remaining))/2))

        // perfom output to screen
        println(status + pad + remaining + pad + timerString)
        if ((time.Duration(timerLength) * time.Second) - elapsed) <= 0 {
            fmt.Print(term.Random(page+line)[:len(page+line)-2])
            if !played {
                f, err := os.Open("notif-sound/notif.mp3")
                if err != nil {
                    log.Fatal(err)
                }
                streamer, format, err := mp3.Decode(f)
                if err != nil {
                    log.Fatal(err)
                }
                defer streamer.Close()

                speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
                speaker.Play(streamer)
                played = true
            }


        } else {
            fmt.Print(page+line)
        }
    }
}
