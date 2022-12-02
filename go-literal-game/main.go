// go-literal-game 一个文字游戏
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Location string
type Forward string
type Chan string
type Object string

const (
	LocationLivingRoom Location = "living-room"
	LocationGarden     Location = "garden"
	LocationAttic      Location = "attic"

	LocationBody Location = "body"

	ForwardWest       Forward = "west"
	ForwardEast       Forward = "east"
	ForwardDownstairs Forward = "downstairs"
	ForwardUpstairs   Forward = "upstairs"

	ChanDoor     Chan = "door"
	ChanStairway Chan = "stairway"

	ObjectWhiskeyBottle Object = "whiskey-bottle"
	ObjectBucket        Object = "bucket"
	ObjectFrog          Object = "frog"
	ObjectChain         Object = "chain"

	ObjectWell   Object = "well"
	ObjectWizard Object = "wizard"
)

type LocationScene struct {
	Description string
	Paths       []*Path
}

type Path struct {
	Forward  Forward
	Chan     Chan
	Location Location
}

var globalMap = map[Location]*LocationScene{
	LocationLivingRoom: {
		Description: "you are in the living room of a wizards house \n " +
			"there is a wizard snoring loudly on the couth \n",
		Paths: []*Path{
			{Forward: ForwardWest, Chan: ChanDoor, Location: LocationGarden},
			{Forward: ForwardUpstairs, Chan: ChanStairway, Location: LocationAttic},
		},
	},
	LocationGarden: {
		Description: "you are in a beautiful garden \n " +
			"there is a well in front of you \n",
		Paths: []*Path{
			{Forward: ForwardEast, Chan: ChanDoor, Location: LocationLivingRoom},
		},
	},
	LocationAttic: {
		Description: "you are in the attic of the wizards house \n " +
			"there is a giant welding torch in the corner \n",
		Paths: []*Path{
			{Forward: ForwardDownstairs, Chan: ChanStairway, Location: LocationLivingRoom},
		},
	},
}

var objectLocations = map[Object]Location{
	ObjectWhiskeyBottle: LocationLivingRoom,
	ObjectBucket:        LocationLivingRoom,
	ObjectFrog:          LocationGarden,
	ObjectChain:         LocationGarden,
}

var chainWelded = false
var bucketFilled = false
var isEnd = false

func look(location Location) {
	locScene, ok := globalMap[location]
	if !ok {
		fmt.Printf("invalid location:%s\n", location)
	}

	fmt.Printf(locScene.Description)
	fmt.Printf("--------------------------------\n")
	for _, path := range locScene.Paths {
		fmt.Printf("there is a %s going 【\033[1;33;40m%s\033[0m】 from here \n", path.Chan, path.Forward)
	}
	fmt.Printf("--------------------------------\n")
	for obj, loc := range objectLocations {
		if location == loc {
			fmt.Printf("you see a 【\u001B[1;33;40m%s\u001B[0m】 on the floor \n", obj)
		}
	}
}

func walkDirection(location Location, forward Forward) Location {
	locScene, ok := globalMap[location]
	if !ok {
		fmt.Printf("invalid location:%s\n", location)
		return location
	}
	var next Location = ""
	for _, path := range locScene.Paths {
		if path.Forward == forward {
			next = path.Location
			break
		}
	}
	if next == "" {
		fmt.Printf("you cannot go that way.\n")
		return location
	}
	return next
}

func pickupObject(location Location, object Object) {
	if objectLocations[object] == location {
		objectLocations[object] = LocationBody
		fmt.Printf("you are now carrying the %s\n", object)
		return
	}
	fmt.Printf("you cannot get that.\n")
}

//func welded(location Location, subject, object Object) {
//	if location == LocationAttic && subject == ObjectChain && object == ObjectBucket &&
//		pickUpObjects[ObjectChain] && pickUpObjects[ObjectBucket] && !chainWelded {
//		chainWelded = true
//		fmt.Printf("the chain is now securely welded to the bucket\n")
//	} else {
//		fmt.Printf("you cannot welded like that\n")
//	}
//}
//
//func dunk(location Location, subject, object Object) {
//	if location == LocationAttic && subject == ObjectBucket && object == ObjectWell &&
//		pickUpObjects[ObjectBucket] && chainWelded {
//		bucketFilled = true
//		fmt.Printf("the bucket is now full of water\n")
//	} else {
//		fmt.Printf("you cannot dunk like that\n")
//	}
//}

type GameActionFunc func(location Location, subject, object Object)

func gameAction(actionName string, place Location, givenSubject, givenObject Object, f func()) func(location Location, subject, object Object) {
	return func(location Location, subject, object Object) {
		if location == place && subject == givenSubject && object == givenObject && objectLocations[givenSubject] == LocationBody {
			f()
		} else {
			fmt.Printf("you cannot %s like that\n", actionName)
		}
	}
}

var gameActionMap = map[string]GameActionFunc{
	"weld":   weld,
	"dunk":   dunk,
	"splash": splash,
}

var weld = gameAction("weld", LocationAttic, ObjectChain, ObjectBucket, func() {
	if objectLocations[ObjectBucket] == LocationBody && !chainWelded {
		chainWelded = true
		fmt.Printf("the chain is now securely welded to the bucket\n")
	} else {
		fmt.Printf("you do not have a bucket or chain is already welded\n")
	}
})

var dunk = gameAction("dunk", LocationGarden, ObjectBucket, ObjectWell, func() {
	if chainWelded {
		bucketFilled = true
		fmt.Printf("the bucket is now full of water\n")
	} else {
		fmt.Printf("the water level is too low to reach\n")
	}
})

var splash = gameAction("splash", LocationLivingRoom, ObjectBucket, ObjectWizard, func() {
	if !bucketFilled {
		fmt.Printf("the bucket has nothing in it\n")
	} else if objectLocations[ObjectFrog] == LocationBody {
		fmt.Printf("the wizard awakens and sees that you stole his frog \n" +
			"he is so upset he banishes you to the netherworlds \n " +
			"you lose! the end \n")
		isEnd = true
	} else {
		fmt.Printf("the wizard awakens from his slumber and greets you warmly \n" +
			"he hands you the magic low-carb donut \n " +
			"you win! the end \n")
		isEnd = true
	}
})

func main() {
	location := LocationLivingRoom
	look(location)
	for !isEnd {
		fmt.Printf("\nPlease input your action [quit look go pick weld dunk splash] : ")
		inputReader := bufio.NewReader(os.Stdin)
		actions, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("invalid input, err:", err)
			continue
		}
		actions = strings.TrimSpace(actions)
		cleanActions := make([]string, 0)
		for _, token := range strings.Split(actions, " ") {
			cleanToken := strings.TrimSpace(token)
			if cleanToken == "" {
				continue
			}
			cleanActions = append(cleanActions, cleanToken)
		}
		if len(cleanActions) == 0 {
			continue
		}
		args := cleanActions[1:]
		switch cleanActions[0] {
		case "quit":
			isEnd = true
		case "look":
			look(location)
		case "go":
			if len(args) < 1 {
				fmt.Println("cannot walk because not given the direction. ")
				continue
			}
			location = walkDirection(location, Forward(args[0]))
			look(location)
		case "pick":
			if len(args) < 1 {
				fmt.Println("cannot pick because not given the object. ")
				continue
			}
			pickupObject(location, Object(args[0]))
		default:
			action := gameActionMap[cleanActions[0]]
			if action == nil {
				fmt.Println("cannot do the action: ", cleanActions[0])
				continue
			}
			if len(args) < 2 {
				fmt.Println("cannot do action because not given the subject and object. ")
				continue
			}
			action(location, Object(args[0]), Object(args[1]))
		}
	}
}
