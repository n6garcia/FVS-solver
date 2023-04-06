package main

func main() {

	//dict := LoadDict()
	dict := LoadWNDict()

	//Solve(dict, "old/")

	//reconstructWord(dict, "happy", "old/delNodes.json")

	exportSol(dict, "wn/delNodes.json", "wnSol.json")

	//simulatedAnnealing(dict, "wn/delNodes.json", "wn/")

	//cullSolution(dict, "old/delNodes.json", "old/")

	//graphVerify(dict, "old/cullNodes.json")

	//alternateVerify(dict, "old/delNodes.json")

	//dictVerify(dict, "wn/cullNodes.json")

	//exportJson(dict)

	//exportCSV(dict, "old/delNodes.json", "old/")

	//handleServer(dict, "old/delNodes.json")
}
