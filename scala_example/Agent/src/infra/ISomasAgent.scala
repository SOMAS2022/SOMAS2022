package infra

import util.Random
import scala.reflect.ClassTag
import otp.OtpHelpers._
import com.ericsson.otp.erlang._

abstract class ISomasAgent(ip: String, cookie: String, serverName: String):
  private val uid = Random.alphanumeric.take(16).mkString
  private val conStr = s"$uid@$ip"
  private val remStr = s"$serverName@$ip"
  private val localNode = new OtpSelf(conStr, cookie)
  private val serverNode = new OtpPeer(remStr)
  private val connection = localNode.connect(serverNode)
  
  private val pid = localNode.pid
  
  def exit() = connection.close()
  Runtime.getRuntime.addShutdownHook(new Thread {override def run = exit()})
  connection.send(serverName, (a"register", pid, localNode.node.eAtm).toOTP)

  val server = SomasServer(connection, connection.receive().asInstanceOf[OtpErlangPid])

  val peers = collection.mutable.Buffer[SomasAgent]()

  def sendMessage(message: SomasMessage, dest: SomasProcess): Unit = {
    dest.conn.send(dest.pid, message.toOtp)
  }

  def run(): Unit
  def receive(msg: SomasMessage): Unit
  
  