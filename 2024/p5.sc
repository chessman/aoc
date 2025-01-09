import scala.util.Sorting
val contestRaw = scala.io.Source
  .fromFile("input.txt")
  .getLines()
  .toList

val input = contestRaw

val rules = input.takeWhile(_ != "").map { s =>
  val p = s.split("\\|").toList.map(_.toInt)
  p(0) -> p(1)
}
val updates =
  input.dropWhile(_ != "").drop(1).map(_.split(",").map(_.toInt).toList)

val rulesMap = rules.foldLeft(Map.empty[Int, List[Int]]) { case (m, (i, j)) =>
  m.get(i) match {
    case Some(l) => m.updated(i, j :: l)
    case None    => m.updated(i, List(j))
  }
}

def isBefore(i: Int, j: Int): Boolean = {
  rulesMap.get(i) match {
    case Some(l) => l.contains(j)
    case None    => false
  }
}

val ord = new Ordering[Int] {
  def compare(x: Int, y: Int): Int = {
    if (isBefore(x, y)) -1
    else if (isBefore(y, x)) 1
    else 0
  }
}

object part1 {
  def run() =
    val ans = updates.map { u =>
      if (u == u.sorted(ord)) u(u.size / 2)
      else 0
    }.sum

    println(ans)
}

object part2 {
  def run() =
    val ans = updates.map { u =>
      val sorted = u.sorted(ord)
      if (u != sorted) sorted(sorted.size / 2)
      else 0
    }.sum

    println(ans)
}

part2.run()


