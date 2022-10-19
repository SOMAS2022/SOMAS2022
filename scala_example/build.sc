import mill._, scalalib._

object Agent extends ScalaModule {
  def scalaVersion = "3.1.1"
  def unmanagedClasspath = T {
    if (!os.exists(millSourcePath / "lib")) Agg()
    else Agg.from(os.list(millSourcePath / "lib").map(PathRef(_)))
  }
}