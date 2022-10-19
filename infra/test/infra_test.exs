defmodule InfraTest do
  use ExUnit.Case
  doctest Infra

  test "greets the world" do
    assert Infra.hello() == :world
  end
end
