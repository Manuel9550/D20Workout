package entities

type Exercise struct {
	RollNumber      int    `json:"rollNumber"`
	ExerciseName    string `json:"exerciseName"`
	StartingAmount  int    `json:"startingAmount"`
	IncrementAmount int    `json:"incrememntAmount"`
	Units           string `json:"units"`
}
