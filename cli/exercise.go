package main

type Exercise struct {
	Id     int    `dynamodbav:"id"`
	Name   string `dynamodbav:"name"`
	Notes  string `dynamodbav:"notes"`
	Sets   int    `dynamodbav:"sets"`
	Reps   []int  `dynamodbav:"reps"`
	Values []int  `dynamodbav:"values"`
	Unit   string `dynamodbav:"unit"`
}
