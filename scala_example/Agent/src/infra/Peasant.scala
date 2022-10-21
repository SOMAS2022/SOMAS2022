package infra

/**
  * This is for infrastructure teams & should be common between agents.
  * These are the types used
  */

sealed trait Stat
case object Hp extends Stat
case object Ap extends Stat
case object At extends Stat
case object Sh extends Stat

sealed trait Boost:
  val stat: Stat

case class PercentBoost(stat: Stat, by: Int) extends Boost
case class FixedBoost(stat: Stat, by: Int) extends Boost

sealed trait Item:
  val boosts: Seq[Boost]

sealed case class Consumable(boosts: Boost*) extends Item

sealed trait Equipment extends Item

sealed abstract class Armour extends Equipment:
  val wear: Int
  val boosts: Seq[Boost]

sealed abstract class Weapon extends Equipment:
  val wear: Int
  val boosts: Seq[Boost]

case class HandArmour(wear: Int, boosts: Boost*) extends Armour
case class BodyArmour(wear: Int, boosts: Boost*) extends Armour
case class HeadArmour(wear: Int, boosts: Boost*) extends Armour
case class LegArmour(wear: Int, boosts: Boost*) extends Armour
case class FootArmour(wear: Int, boosts: Boost*) extends Armour
case class OneHandedWeapon(wear: Int, boosts: Boost*) extends Weapon
case class TwoHandedWeapon(wear: Int, boosts: Boost*) extends Weapon

class Hands:
  val weapon: Either[(Option[OneHandedWeapon], Option[OneHandedWeapon]), Option[TwoHandedWeapon]] = Right(None)

case class Peasant(baseHp: Int, baseAp: Int, baseAt: Int, baseSh: Int, level: Int, items: Item*):
  def getBaseStat(stat: Stat): Int = stat match
    case Hp => baseHp
    case Ap => baseAp
    case At => baseAt
    case Sh => baseSh

  def getStat(stat: Stat): Int =
    getBaseStat(stat) + items.map(_.boosts.collect {
      case PercentBoost(`stat`, by) => (getBaseStat(stat) * by) / 100
      case FixedBoost(`stat`, by) => by
    }.sum).sum