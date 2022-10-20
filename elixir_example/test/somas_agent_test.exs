defmodule SomasAgentTest do
  use ExUnit.Case
  doctest SomasAgent

  test "greets the world" do
    assert SomasAgent.hello() == :world
  end
end
