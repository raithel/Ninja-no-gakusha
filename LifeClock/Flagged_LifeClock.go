/*
 Copyright 2016 Christopher M. Kruczek <krucz7@gmail.com>

  This software is being distributed under the:
	GNU General Public License v3.
	For full details of the license.
	see: http://www.gnu.org/copyleft/gpl.html
 Date : 1-21-2016

 Program Description :  A Clock that calculates how long you have been alive for in days.
							uses flags as a way to input a Birth date.
							can use a .sh or .bat file.
 golang as of version: 1.5.3

TODO:
	[Done] Rework LifeClock.go to work with Flags.
	[Done] Use escape characters: http://stackoverflow.com/questions/15442292/golang-how-to-have-an-inplace-string-that-updates-at-stdout
  - enable cleaning/clearing of the screen periodically to fix overlap of characters.
	- Implement a cross platform version that uses .sh and .bat files.
		- Implement window resizing via .sh file.
		 
FUNCTION:
	- Program takes a date and time as flags for input.
	- Inputted time is used for calculating a persons age in seconds.
  -
	-
  -


*/
//package io
//has method writer
// can be used with fmt.Fprintf()

package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

var bigNums = [5][11]string{{" ___ ", "     ", " ___ ", " ___ ", "     ", " ___ ", " ___ ", " ____", " ___ ", " ___ ", "     "},
	{"/  /\\", " /|  ", "    \\", "    \\", "|   |", "/    ", "/    ", "    /", "/   \\", "/   |", "  o  "},
	{"| / |", "  |  ", " ___|", " ___|", "|___|", "|___ ", "|___ ", "   / ", "\\___/", "\\___|", "     "},
	{"|/  |", "  |  ", "|    ", "    |", "    |", "    |", "|   |", "  /  ", "/   \\", "    |", "  o  "},
	{"\\___/", "__|__", "\\___ ", " ___/", "    |", " ___/", "\\___/", " /   ", "\\___/", "    |", "     "}}

func checkTime(now *int, later *int) bool { // Compares two times for differences.
	if *now != *later {
		return true
	}
	return false
}

func waitForChange(curSec *int) bool { // Makes the program wait until the time changes.
	var (
		timeChange = false
		secNow     int
	)
	for timeChange == false { // Loops until time changes.
		secNow = time.Now().Second()            // Gets current time in seconds.
		time.Sleep(50 * 10000000)               // Waits for one second in microseconds.
		timeChange = checkTime(curSec, &secNow) // Compares two times, for a difference.
	}
	return timeChange
}

func getDurNums(dur *int) (*int, *int, *int) { // [Get Duration Numbers] Takes a time in seconds and returns hours minutes and seconds.
	var (
		hour   = *dur/60/60%24 - 4 // int
		minute = *dur / 60 % 60    // int
		sec    = *dur % 60         // int
	)
	return &hour, &minute, &sec
}

func getTimeNums(cTime *time.Time) (*int, *int, *int) { // Takes the given time and returns hours, minutes and seconds.
	var (
		hour   = cTime.Hour()   //int
		minute = cTime.Minute() //int
		sec    = cTime.Second() //int
	)
	if hour > 12 { // adjusts military time to normal time.
		hour -= 12
	}
	return &hour, &minute, &sec
}

func formatNumber(num *int) string { // Converts an int to a string.
	if *num < 10 {
		return "0" + strconv.Itoa(*num)
	}
	return strconv.Itoa(*num)
}

func printTime(hour *int, minute *int, sec *int) { // Takes units of time, and prints them as a string.
	fmt.Printf("%s:%s:%s", strconv.Itoa(*hour), formatNumber(minute), formatNumber(sec))
}
func printBigTime(hour *int, minute *int, sec *int) { // Displays the time in large ascii art numbers.
	var (
		tinyClock = "" + strconv.Itoa(*hour) + ":" + formatNumber(minute) + ":" + formatNumber(sec) // string
		index     int                                                                               // The number that is used to access the Bignums array.
		output    [5]string                                                                         // An array of five strings used to output the large asciiart time.
	)
	for _, value := range tinyClock { // Parses through tinyClock string and adds the appropriate number/symbol to the output array.
		if string(value) == ":" { // If the character is a colon, then set the index to the colon ascii art to the output array.
			index = 10
		} else {
			index, _ = strconv.Atoi(string(value)) // Parses string to an int.
		}
		for i, line := range bigNums { // Appends the string to the output array of strings.
			//fmt.Printf("%s\n",line[index])
			output[i] += (line[index])
		}
	}
	for _, line := range output { // Displays the output array of strings.
		fmt.Printf("%s\n", line)
	}
	fmt.Printf(tinyClock) // Displays the time in a small standard letter format.
}


func clear() { // Moves the cursor to the point 0,0 on the terminal, using an escape sequence.
	fmt.Printf("\033[0;0H") // row;columnH
}

func main() {
	var (
		day     int                                                                     // Used for flag -d
		month   int                                                                     // Used for flag -m
		year    int                                                                     // Used for flag -y
		hour    int                                                                     // Used for flag -hr
		minute  int                                                                     // Used for flag -min
		banner  string                                                                  // Used for flag -banner
		t       = time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC) // time.Time
		curTime = time.Now()                                                            // time.Time
		sec     int                                                                     // Used to compare current time with previous time.
	)
	// Setup flags.
	flag.IntVar(&day, "d", 01, "day of Birth")
	flag.IntVar(&year, "y", 1992, "year of Birth")
	flag.IntVar(&month, "m", 06, "month of Birth")
	flag.IntVar(&hour, "hr", 06, "hour of Birth")
	flag.IntVar(&minute, "min", 00, "minute of Birth")
	flag.StringVar(&banner, "banner", "Your Life Clock", "Displayed Banner")

	flag.Parse() // Parse inputted flags.

	// function inside of a function.
	//  This needed to be done, to allow for the date [t] inputted by flags to be used in the main function.
	getSeconds := func(curTime *time.Time) *int { // Takes the current time and returns the Birth altered time in seconds.
		var d = int(curTime.Sub(t).Seconds()) // Gets the current time minus the birth time in seconds.
		fmt.Printf("%d Days\n", d/86400)      // Displays how many days you have been alive for. 86400 is the equivalent of /60/60/24.
		return &d
	}

	for {
		sec = curTime.Second()                         // get current time in seconds.
		waitForChange(&sec)                            // loop until there is a change in seconds.
		curTime = time.Now()                           // When there is a change in seconds, update the time.
		clear()                                        // Return the cursor to the uppermost right hand side of the terminal.
		fmt.Printf("%s:\n", banner)                    // Display Banner.
		printBigTime(getDurNums(getSeconds(&curTime))) // display the big time.
		fmt.Printf("\nStandard Time ")
		printTime(getTimeNums(&curTime))
		//fmt.Println(t)//enable for debugging
	} // end infinite for loop.
} // end main.
