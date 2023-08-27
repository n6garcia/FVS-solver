package main

func main() {

	//handleServer("wnSol.json")

	dict := LoadDict()
	//dict := LoadWNDict()

	Solve(dict)

	//reconstructWord(dict, "happy", "delNodes.json")

	//exportSol(dict, "delNodes.json", "oldSol.json")

	//simulatedAnnealing(dict, "delNodes.json")

	//cullSolution(dict, "delNodes.json")

	graphVerify(dict, "delNodes.json")

	//alternateVerify(dict, "delNodes.json")

	//dictVerify(dict, "cullNodes.json")

	//exportTrees(dict, "delNodes.json")

	//exportNames(dict)

	//exportJson(dict)

	//exportCSV(dict, "delNodes.json")
}
