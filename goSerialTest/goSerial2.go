package main

import (
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "sync"

    "github.com/tarm/serial"
)

const DEFAULT_BAUD_RATE = 9600

func main() {
    flag.Usage = func() {
        name := filepath.Base(os.Args[0])
        fmt.Fprintf(os.Stderr, "USAGE: %s <serial-device> [options]\n", name)
        fmt.Fprint(os.Stderr, "Options:\n")
        flag.PrintDefaults()
        fmt.Fprintln(os.Stderr, "\n")
        fmt.Fprintln(os.Stderr, "Example: ")
        fmt.Fprintf(os.Stderr, "\t %s /dev/ttyUSB0 --baud-rate 115200 \n", name)
    }
    baudRate := flag.Int("baud-rate", DEFAULT_BAUD_RATE, "The speed of the serial device")

    //if len(os.Args) < 2 {
    //    flag.Usage()
    //    os.Exit(1)
    //}

    //devicePath := os.Args[1]
    devicePath := "/dev/ttyS0"

    config := &serial.Config{Name: devicePath, Baud: *baudRate}
    serialPort, err := serial.OpenPort(config)

    if err != nil {
        log.Fatalf("Could not open serial device %s", devicePath)
    }

    defer serialPort.Close()

    wg := sync.WaitGroup{}
    wg.Add(2)

    // redirect the stdin to the serial device
    go func() {
        defer wg.Done()

        if _, err := io.Copy(serialPort, os.Stdin); err != nil {
            log.Fatal(err)
        }
    }()

    // redirect the serial device to stdout
    go func() {
        defer wg.Done()

        if _, err := io.Copy(os.Stdout, serialPort); err != nil {
            log.Fatal(err)
        }
    }()

    wg.Wait()
}