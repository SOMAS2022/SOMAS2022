import com.ericsson.otp.erlang._
import infra.otp.Node

@main def main =
  val node = Node("127.0.0.1", "QKPYDTHMTBVMEDRPLKQD", "infra_server")
  SomasAgent(node).run()
  node.connection.close()