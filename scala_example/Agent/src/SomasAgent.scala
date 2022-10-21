import infra._

class SomasAgent(ip: String, cookie: String, serverName: String = "infra_server") extends ISomasAgent(ip, cookie, serverName):
  // TODO: Implement this :)
  def run(): Unit = ???
  def receive(message: SomasMessage): Unit = ???