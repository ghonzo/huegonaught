package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	hue "github.com/collinux/gohue"
)

type SwitchState int

const (
	On SwitchState = iota
	TurningOff
	Off
)

func main() {
	bridgeFlag := flag.String("bridge", "philips-hue", "the ip address of the Hue bridge")
	userFlag := flag.String("user", "", "(required) the user for Hue bridge")
	onPollingIntervalFlag := flag.Int("on", 10, "the polling interval (in seconds) when the lights are on")
	offPollingIntervalFlag := flag.Int("off", 60, "the polling interval (in seconds) when the lights are off")
	signalBulbsFlag := flag.String("signal", "", "(required) comma-separated list of bulb ids to monitor for reachable status")
	controlledBulbsFlag := flag.String("controlled", "", "(required) comma-separated list of light ids to turn off when signal bulbs are unreachable")
	verboseFlag := flag.Bool("v", false, "verbose")
	flag.Parse()
	// All of these are required fields
	if len(*userFlag) == 0 || len(*signalBulbsFlag) == 0 || len(*controlledBulbsFlag) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	// Now make sure both bulb id lists are ints
	var err error
	signalBulbsStrs := strings.Split(*signalBulbsFlag, ",")
	signalBulbs := make([]int, len(signalBulbsStrs))
	for i, str := range signalBulbsStrs {
		if signalBulbs[i], err = strconv.Atoi(str); err != nil {
			log.Fatal("Invalid signal bulb id:", str)
		}
	}
	controlledBulbsStrs := strings.Split(*controlledBulbsFlag, ",")
	controlledBulbs := make([]int, len(controlledBulbsStrs))
	for i, str := range controlledBulbsStrs {
		if controlledBulbs[i], err = strconv.Atoi(str); err != nil {
			log.Fatal("Invalid controlled bulb id:", str)
		}
	}
	// Finally set the polling intervals
	onPollingInterval := time.Duration(*onPollingIntervalFlag) * time.Second
	offPollingInterval := time.Duration(*offPollingIntervalFlag) * time.Second
	// Let's go!
	bridge, err := hue.NewBridge(*bridgeFlag)
	if err != nil {
		log.Fatal("Could not connect to bridge", *bridgeFlag, err)
	}
	err = bridge.Login(*userFlag)
	if err != nil {
		log.Fatal("Could not login with given user", err)
	}
	ticker := time.NewTicker(onPollingInterval)
	log.Println("Connected to bridge, starting to poll now")
	switchState := On
	for ; true; <-ticker.C {
		if *verboseFlag {
			log.Println("Polling ... switchState = ", switchState)
		}
		if lights, err := bridge.GetAllLights(); err != nil {
			log.Println("Error getting lights:", err)
		} else {
			lightsMap := make(map[int]hue.Light)
			for _, light := range lights {
				lightsMap[light.Index] = light
			}
			allUnreachable := areAllUnreachable(lightsMap, signalBulbs)
			// See if we flipped from reachable to unreachable
			switch switchState {
			case On:
				if allUnreachable {
					log.Println("Detected switch off, so transitioning to TurningOff")
					turnOffBulbs(lightsMap, controlledBulbs)
					switchState = TurningOff
				}
			case TurningOff:
				// Make sure they are still off
				if !allUnreachable {
					log.Println("Switch flipped back on, so go back to On state")
					switchState = On
				} else if areAllOff(lightsMap, controlledBulbs) {
					log.Println("Confirmed all controlled lights are off, transitioning to Off state with polling interval", offPollingInterval)
					ticker.Reset(offPollingInterval)
					switchState = Off
				} else {
					log.Println("Some lights are still on, sending off signal again")
					turnOffBulbs(lightsMap, controlledBulbs)
					// Stay in this state
				}
			case Off:
				if !allUnreachable {
					log.Println("Detected switch on, so switching polling interval to", onPollingInterval)
					ticker.Reset(onPollingInterval)
					switchState = On
				}
			}
		}
	}
}

// Return false if at least one is reachable, otherwise return true
func areAllUnreachable(lightsMap map[int]hue.Light, signalBulbs []int) bool {
	for _, bulb := range signalBulbs {
		if light, ok := lightsMap[bulb]; ok && light.State.Reachable {
			return false
		}
	}
	return true
}

// Return false if at least one is on, otherwise return true
func areAllOff(lightsMap map[int]hue.Light, signalBulbs []int) bool {
	for _, bulb := range signalBulbs {
		if light, ok := lightsMap[bulb]; ok && light.State.On {
			return false
		}
	}
	return true
}

// Turn off all of the bulbs in the controlledBulbs slice
func turnOffBulbs(lightsMap map[int]hue.Light, controlledBulbs []int) {
	for _, bulb := range controlledBulbs {
		if light, ok := lightsMap[bulb]; !ok {
			log.Printf("Did not find bulb #%d to turn it off\n", bulb)
		} else {
			if err := light.Off(); err != nil {
				log.Printf("Failed to turn off bulb #%d: %v\n", bulb, err)
			}
		}
	}
}
