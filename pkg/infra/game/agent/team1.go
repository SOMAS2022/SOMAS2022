package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/teams/team1"
	"math"
	"math/rand"
	"time"
)

// Function which defines how an agent perceives an action
func actionSentiment(action decision.FightAction) [4]float64 {
	switch action {
	case decision.Cower:
		return [4]float64{0.0, 0.0, -0.1, -0.1}
	case decision.Attack:
		return [4]float64{0.0, 0.0, 0.1, 0.1}
	case decision.Defend:
		return [4]float64{0.0, 0.0, 0.1, 0.1}
	default:
		return [4]float64{0.0, 0.0, 0.0, 0.0}
	}
}

// Ensures a float is between -1 and 1
func boundFloat(inputNumber float64) float64 {
	if inputNumber > 1.0 {
		return 1.0
	} else if inputNumber < -1.0 {
		return -1.0
	} else {
		return inputNumber
	}
}

// Ensures array values are between -1 and 1
func boundArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		boundFloat(inputArray[0]),
		boundFloat(inputArray[1]),
		boundFloat(inputArray[2]),
		boundFloat(inputArray[3]),
	}
}

// Add two arrays
func addArrays(A [4]float64, B [4]float64) [4]float64 {
	return [4]float64{
		A[0] + B[0],
		A[1] + B[1],
		A[2] + B[2],
		A[3] + B[3],
	}
}

func decayNumber(inputNumber float64) float64 {
	if inputNumber < 0 {
		return 0.70 * inputNumber
	} else {
		return 0.90 * inputNumber
	}
}

func decayArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		decayNumber(inputArray[0]),
		decayNumber(inputArray[1]),
		decayNumber(inputArray[2]),
		decayNumber(inputArray[3]),
	}
}

func softmax(inputArray [3]float64) [3]float64 {
	expValues := [3]float64{
		math.Exp(inputArray[0]),
		math.Exp(inputArray[1]),
		math.Exp(inputArray[2]),
	}

	// Sum exponential array
	sum := 0.0
	for i := 0; i < 3; i++ {
		sum += expValues[i]
	}

	// Divide each element in input array by sum
	for i := 0; i < 3; i++ {
		expValues[i] /= sum
	}

	return expValues
}

func makeIncremental(inputArray [3]float64) [3]float64 {

	var outputArray [3]float64

	outputArray[0] = inputArray[0]

	for i := 1; i < 3; i++ {
		outputArray[i] = outputArray[i-1] + inputArray[i]
	}

	return outputArray
}

type Team1Agent struct {
	socialCapital    map[string][4]float64
	lastLevelUpdated uint       // The level at which the socialCapital was last updated
	pCooperation     float64    // Probability that agent will cooperate
	coopTable        [3]float64 // QTable for actions when agent cooperates towards the common goal
	selfTable        [3]float64 // QTable for actions when the agent acts in its own self-interest

	// Four metrics for each agent's perception of other agents. Three metrics are borrowed from Ostrom-Ahn
	// social capital model: (The ordering below is the same as the ordering in the array)
	// 1. Institutions: Influenced by institutional actions such as voting. Also influenced by any elected
	//positions the agent currently holds.
	// 2. Networks: Messages received about an agent from other agents.
	// The impact which a message from another agent has on the network metric depends on the overall perception of
	// the agent sending the message.
	// 3. Trustworthiness: Measures how likely the agent is to cooperate based on historic decisions.
	// Finally, our addition to the model is
	// 4. Honour: A metric of how many favours an agent has done for another agent, compared to the amount of
	//favours the other agent has done for them. This metric is supposed to capture social expectations, such
	//as if I went out with the trash last time, then I feel that next time someone else should do it. The value
	// of the  favours metric will over time impact the trustworthiness metric of another agent. In many ways,
	// favours is a short time horizon equivalent of trustworthiness.

	// Forgiveness: In order to enable even the bitterest of enemies to start cooperating again in the long term,
	// the metrics will decay with time. However, negative metrics will decay far faster than positive metrics.
	// This is similar to how you can have a massive fight and hate someone one day, and then be friends again the
	// next. On the other hand, you can go years without seeing a friend and still have a very positive impression
	// of them.

	// Messages we should have: Denounce and Praise
}

// Calculate utility value of different decisions
func (r Team1Agent) utilityValue(action decision.FightAction, view *state.View, agent BaseAgent) float64 {
	// Utility of each action is dependent on relationship with others. If agent hates all other agents, then
	// will only act in its own interest.

	agentStats, _ := view.AgentState().Get(agent.Id)

	switch action {
	case decision.Cower:
		// Goes down with health, and down with stamina and down with high social capital of others
		return 0.005 * float64(1000-int(agentStats.Hp))
	case decision.Attack:
		// Goes up with health, and up with stamina and up with high social capital of others
		return 0.005 * float64(int(agentStats.Hp))
	case decision.Defend:
		return 0.005 * float64(int(agentStats.Hp))
	default:
		return 1
	}
}

// Called any time a message is received, initialises or updates the socialCapital map
func (r Team1Agent) updateSocialCapital(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {

	// Ensure that socialCapital map is initialised
	if len(r.socialCapital) == 0 && view.AgentState().Len() > 1 {

		// Create empty map
		r.socialCapital = map[string][4]float64{}

		// Populate map with every currently living agent, and calculate socialCapital based on log
		itr := view.AgentState().Iterator()
		for !itr.Done() {
			key, _, _ := itr.Next()

			action, exists := log.Get(key)
			if exists { // If agent exists in log, calculate socialCapital
				r.socialCapital[key] = boundArray(actionSentiment(action))
			} else { // Else initialize socialCapital to 0
				r.socialCapital[key] = [4]float64{0.0, 0.0, 0.0, 0.0}
			}
		}

		// Delete the agents own id from the socialCapital array
		delete(r.socialCapital, agent.Id)

		// Set the lastLevelUpdated variable
		r.lastLevelUpdated = view.CurrentLevel()

	} else if r.lastLevelUpdated < view.CurrentLevel() { // socialCapital variable already exists

		for key := range r.socialCapital {

			// Remove any agents that have died from socialCapital map (Might be unnecessary as it adds a lot of computation)
			_, exists := view.AgentState().Get(key)
			if !exists {
				delete(r.socialCapital, key)
			}

			// Decay socialCapital values
			r.socialCapital[key] = decayArray(r.socialCapital[key])

			// TODO: Update of socialCaptial should be dependent on the agents own action (especially for favours)
			// Update socialCapital based on log
			action, exists := log.Get(key)
			if exists {
				r.socialCapital[key] = addArrays(r.socialCapital[key], boundArray(actionSentiment(action)))
			}
		}

		// Set lastLevelUpdated to current level
		r.lastLevelUpdated = view.CurrentLevel()
	}

	// Ensure all socialCapital values are between -1 and 1
	for key := range r.socialCapital {
		r.socialCapital[key] = boundArray(r.socialCapital[key])
	}

}

// Called whenenver an unimplemented message is received
func (r Team1Agent) Default() decision.FightAction {
	//TODO implement me
	panic("implement me")
}

func NewTeam1Agent() *Team1Agent {

	// Initialise a random seed
	rand.Seed(time.Now().UnixNano())

	// Check if environment variables exist for coopTable and selfTable
	//coopTable        QTable  // QTable for actions when agent cooperates towards the common goal
	//selfTable        QTable  // QTable for actions when the agent acts in its own self-interest
	return &Team1Agent{
		pCooperation: rand.Float64(),
		coopTable:    [3]float64{0.0, 0.0, 0.0},
		selfTable:    [3]float64{0.0, 0.0, 0.0},
	}
}

func (r Team1Agent) HandleFightMessage(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) decision.FightAction {

	r.updateSocialCapital(m, view, agent, log)

	// Decide if agent is going to cooperate or act in own self-interest
	cooperate := rand.Float64() < r.pCooperation

	var utilityValues [3]float64
	if cooperate {
		utilityValues = r.coopTable
	} else {
		utilityValues = r.selfTable
	}

	// Apply softmax to get probabilities
	softArray := softmax(utilityValues)

	// Make number representation incremental
	probArray := makeIncremental(softArray)

	/*it := view.AgentState().Iterator()
	nextId, _, _ := it.Next()
	if agent.Id == nextId {
		fmt.Println(utilCower)
		fmt.Println([3]float64{utilCower, utilAttack, utilDefend})
		fmt.Println(softArray)
		fmt.Println(probArray)
	}*/

	// Do action with probability based on utility value
	var action decision.FightAction
	switch random := rand.Float64(); {
	case 0.0 < random && random < probArray[0]:
		action = decision.Cower
	case probArray[0] < random && random < probArray[1]:
		action = decision.Attack
	case probArray[1] < random && random < probArray[2]:
		action = decision.Defend
	}

	// Maybe only log with a probability, otherwise logs become very long
	team1.LogAction(view, agent.Id, action)
	//team1.LogAction()

	return action
}
