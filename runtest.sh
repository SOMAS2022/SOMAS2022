# each el in list are params for a game
listAgentCollective=(10 0 20)
listAgentSelfish=(20 10 0)
listAgentAgr=(0 20 10)

listSELFISH_PER=(25 25 25)
listCOLLECTIVE_PER=(50 50 50)
listSELFLESS_PER=(75 75 75)

listDEFECTION=(true true true)
listUPDATE_PERSONALITY=(true true true)

for i in ${!listAgentCollective[@]}; do
  AgentCollective=${listAgentCollective[i]}
  AgentSelfish=${listAgentSelfish[i]}
  AgentAgr=${listAgentAgr[i]}

  AgentSELFISH_PER=${listSELFISH_PER[i]}
  AgentCOLLECTIVE_PER=${listCOLLECTIVE_PER[i]}
  AgentSELFLESS_PER=${listSELFLESS_PER[i]}

  AgentDEFECTION=${listDEFECTION[i]}
  AgentUPDATE_PERSONALITY=${listUPDATE_PERSONALITY[i]}
  
  export AGENT_TEAM3NEUT_QUANTITY=$AgentCollective
  export AGENT_TEAM3PAS_QUANTITY=$AgentAgr
  export AGENT_TEAM3AGR_QUANTITY=$AgentSelfish
 
  export SELFISH_PER=$AgentSELFISH_PER
  export COLLECTIVE_PER=$AgentCOLLECTIVE_PER
  export SELFLESS_PER=$AgentSELFLESS_PER

  export DEFECTION=$AgentDEFECTION
  export UPDATE_PERSONALITY=$AgentUPDATE_PERSONALITY

  OUTPUT=$(make| tail -1)
  echo $OUTPUT
  python3 plotGame.py -f $OUTPUT

done
