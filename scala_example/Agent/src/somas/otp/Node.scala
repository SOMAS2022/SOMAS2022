package somas.otp
import util.Random
import com.ericsson.otp.erlang.OtpSelf
import com.ericsson.otp.erlang.OtpPeer
import scala.reflect.ClassTag
import OtpHelpers._
import com.ericsson.otp.erlang.OtpErlangTuple

class Node(ip: String, cookie: String, serverName: String = "infra_server"):
  val uid = Random.alphanumeric.take(16).mkString
  val conStr = s"$uid@$ip"
  val localNode = new OtpSelf(conStr, cookie)
  println(s"Set up at $conStr")
  val remStr = s"$serverName@$ip"
  val serverNode = new OtpPeer(remStr)
  println(s"Connecting to $remStr")

  val connection = localNode.connect(serverNode)
  
  val pid = localNode.pid

  def hasMsg() = connection.msgCount() > 0
  def isAlive() = connection.isAlive()
  def getTupleMsg() = connection.waitFor[OtpErlangTuple].toList
  def send(msg: Any*) = {
    val otpMsg = msg.toOTP
    connection.send(serverName, otpMsg)
  }
  def exit() = connection.close

  Runtime.getRuntime.addShutdownHook(new Thread {override def run = exit()})
  send(a"register", pid, localNode.node)