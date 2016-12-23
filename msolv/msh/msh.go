package msh

// configuration specific to MathSolver Marsha

// MathProblem is the type of message sent to the MathProblemsA bulletin board.
type MathProblem struct {
	Name, ID string
	Data     []byte
}

//MathAnswer is an answer to a MathProblem
type MathAnswer struct {
	SolverID  string
	ProblemID string
	AnswerID  string
	Answer    []byte
}
