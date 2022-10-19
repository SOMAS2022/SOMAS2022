defmodule Agent do
  @moduledoc """
    Contains agent-specific implementations
  """

  @typedoc """
    Agent struct - contains info pertaining to the connection and status of each agent currently in the dungeon
  """
  defstruct [:node, :hp, :ap]

  @doc """
    Start the current round, and wait for a response from the agent
  """
  def start_round(agent) do

  end

end
