package infra

trait ISomasAgent(node: otp.Node):
  def run(): Unit
  