package main

import "go-lab/app/go_wire/internal/wiref"

func main() {
	mission := wiref.InitMissionWithValue("lbb")
	mission.Start()
}
