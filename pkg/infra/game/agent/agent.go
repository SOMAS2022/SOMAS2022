package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/game/strategy"
	"sync"

	"github.com/benbjohnson/immutable"
)

type Agent struct {
	BaseAgent message.BaseAgent
	Strategy  strategy.Strategy
}

func (a *Agent) HandleFight(
	agentMap map[commons.ID]Agent,
	view state.View,
	log immutable.Map[commons.ID, decision.FightAction],
	decisionChan chan strategy.ActionDecision,
	wg *sync.WaitGroup) {
	// Process any messages that the agent currently has in a loop
	// TODO ? keep looping this for loop until all agents are completed
	for m := range a.BaseAgent.Communication.Receipt {
		a.handleMessage(agentMap, &view, &log, m)
		action := a.Strategy.GenerateActionDecision()
		if action != decision.Undecided {
			go func() {
				// TODO add switch here on message type
				<-a.BaseAgent.Communication.Receipt
			}()
			decisionChan <- strategy.ActionDecision{Action: action, Sender: a.BaseAgent.Id}
			wg.Done()
			return
		}
	}

	// Ask the agent for its final decision
	decisionChan <- strategy.ActionDecision{Action: a.Strategy.GenerateActionDecision(), Sender: a.BaseAgent.Id}
}

func (a *Agent) handleMessage(agentMap map[commons.ID]Agent, view *state.View, log *immutable.Map[commons.ID, decision.FightAction], m message.TaggedMessage) {
	switch m.Message.(type) {
	case strategy.RequestMessageInterface:
		// TODO add timeout in which agent must reply
		response := m.Message.(strategy.RequestMessageInterface).ProcessRequestMessage(a.Strategy, view, log)
		response.SetUUID(m.Message.GetUUID())
		//todo: handle error in case of message failure
		agent := agentMap[m.Sender].BaseAgent
		_ = message.SendBlockingMessageImp(&agent, a.BaseAgent.Id, response)

	case strategy.InfoMessageInterface:
		m.Message.(strategy.InfoMessageInterface).ProcessInfoMessage(a.Strategy, view, log)
	default:
	}
}
