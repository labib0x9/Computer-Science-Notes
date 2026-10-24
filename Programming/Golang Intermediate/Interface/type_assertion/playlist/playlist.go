package main 

type TapePlayer struct {

}

type TapeRecorder struct {

}

type Playable interface {
	Play(song string)
	Stop()
}

func main() {

	songs := []string { "song 1", "song 2" }
	
	tapePlayer := TapePlayer{}
	PlayList(tapePlayer, songs)

	rsongs := []string { "song 4", "song 3" }
	recorder := TapeRecorder{}
	PlayList(recorder, rsongs)
}

func PlayList(device Playable, songs []string) {

	for _, song := range songs {
		device.Play(song)
		
		// Type Assertion
		if recorder, ok := device.(TapeRecorder); ok {
			recorder.Record()
		}
	}

	device.Stop()
}

func (t TapePlayer) Play(song string) {

}

func (t TapePlayer) Stop() {

}

func (t TapeRecorder) Play(song string) {

}

func (t TapeRecorder) Stop() {

}

func (t TapeRecorder) Record() {

}