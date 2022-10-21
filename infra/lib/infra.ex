defmodule Infra do
  @moduledoc """
  Documentation for `Infra`.
  """
  use Application

  defstruct [:n_agents, :agents, :dungeon, :time]

  @doc """
  Loop function - handles requests from
  """
  @spec loop(state::Infra.struct()) :: no_return()
  def loop(state) do
    receive do
      {:register, process, node} ->
        :io.fwrite("Node: ~p (atomic? ~p), Process; ~p~n", [node, is_atom(node), process])
        # TODO: Put initial state in here
        #send(process, {:welcome, :agent})
        send(process, self())
        pid = :rpc.call(node, SomasAgent, :init, [])
        :io.fwrite("Pid is: ~p~n", [pid])
      msg ->
        :io.fwrite("Received ~p~n", [msg])
    end
    loop(state)
  end

  @doc """
    Initial function - registers the infra server globally so it can be called by agents.
  """
  @spec init(list()) :: no_return()
  def init(_args\\[]) do
    :global.register_name(:infra_server, self())
    :erlang.register(:infra_server, self())
    initial_state = %Infra{}
    loop(initial_state)
  end

  def start(_type, args) do
    # TODO: Resolve args to customise # of agents, etc.
    {:ok, spawn(fn -> init(args) end)}
  end

  def stop(_state) do
    :ok
  end
end
