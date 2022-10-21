package infra
import com.ericsson.otp.erlang.OtpErlangPid
import com.ericsson.otp.erlang.OtpErlangObject
import com.ericsson.otp.erlang.OtpConnection

sealed trait SomasProcess:
  val conn: OtpConnection
  val pid: OtpErlangPid

case class SomasServer(conn: OtpConnection, pid: OtpErlangPid) extends SomasProcess

case class SomasAgent(conn: OtpConnection, pid: OtpErlangPid) extends SomasProcess

sealed trait SomasMessage:
  def toOtp: OtpErlangObject

sealed trait SomasServerMessage extends SomasMessage
sealed trait SomasAgentMessage extends SomasMessage