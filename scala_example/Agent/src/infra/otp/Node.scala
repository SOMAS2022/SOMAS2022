package infra.otp
import util.Random
import scala.reflect.ClassTag
import OtpHelpers._
import com.ericsson.otp.erlang._

class Node(ip: String, cookie: String, serverName: String = "infra_server"):
  val uid = Random.alphanumeric.take(16).mkString
  val conStr = s"$uid@$ip"
  val localNode = new OtpSelf(conStr, cookie)
  localNode.publishPort()
  println(s"Set up at $conStr")
  val remStr = s"$serverName@$ip"
  val serverNode = new OtpPeer(remStr)
  println(s"Connecting to $remStr")

  val connection = localNode.connect(serverNode)

  connection.sendRPC("io", "fwrite", Array(OtpErlangBinary("Hello World~n").asInstanceOf[OtpErlangObject]))
  
  val pid = localNode.pid

  def hasMsg() = connection.msgCount() > 0
  def isAlive() = connection.isAlive()
  def send(msg: Any*) = {
    val otpMsg = msg.toOTP
    connection.send(serverName, otpMsg)
  }
  def exit() = connection.close

  Runtime.getRuntime.addShutdownHook(new Thread {override def run = exit()})
  send(a"register", pid, localNode.node.eAtm)