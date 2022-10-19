import com.ericsson.otp.erlang._
import somas.otp.Node

@main def main =
  val node = new Node("127.0.0.1", "QKPYDTHMTBVMEDRPLKQD", "infra_server")
  println(node.isAlive())
  println(node.uid)
  val recv = node.getTupleMsg()
  println(recv)
