defmodule Infra do
  @moduledoc """
  Documentation for `Infra`.
  """

  @doc """
    Loop function - handles requests from
  """
  def loop(state) do
    receive do
      {:register, process, node} ->

      msg ->
        :io.fwrite("Received ~p~n", [msg])
    end
    loop(state)
  end

  @doc """
    Initial function - registers the infra server globally so it can be called by agents.
  """
  def init do
    :global.register_name(:infra_server, self())
    :erlang.register(:infra_server, self())
    initial_state = %{}
    loop(initial_state)
  end
end
