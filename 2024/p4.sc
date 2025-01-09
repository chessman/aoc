object v1 {}

val demo = """
MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX""".split("\n").drop(1).map(_.split("").map(_.head))

val contest = scala.io.Source
  .fromFile("input.txt")
  .mkString
  .split("\n")
  .map(_.split("").map(_.head))

val input = contest

val MaxX = input(0).size - 1
val MaxY = input.size - 1

case class P(x: Int, y: Int)

object part1 {
  def directions = List(
    (p: P) => P(p.x - 1, p.y - 1),
    (p: P) => P(p.x - 1, p.y),
    (p: P) => P(p.x - 1, p.y + 1),
    (p: P) => P(p.x, p.y - 1),
    (p: P) => P(p.x, p.y + 1),
    (p: P) => P(p.x + 1, p.y - 1),
    (p: P) => P(p.x + 1, p.y),
    (p: P) => P(p.x + 1, p.y + 1)
  )
  
  case class S(p: P, dir: P => P)

  def filterByChar(ps: List[P], c: Char): List[P] =
    ps.filter(p => input(p.y)(p.x) == c)

  def advance(ss: List[S], c: Char): List[S] =
    ss.flatMap { s =>
      val newP = s.dir(s.p)
      if (
        newP.x >= 0 &&
        newP.x <= MaxX &&
        newP.y >= 0 &&
        newP.y <= MaxY &&
        input(newP.y)(newP.x) == c
      ) {
        List(S(newP, s.dir))
      } else Nil
    }

  def run() = {
    val allS = (
      for {
        x <- 0 to MaxX
        y <- 0 to MaxY
        dir <- directions
      } yield S(P(x, y), dir)
    ).toList

    val xs = allS.filter(s => input(s.p.y)(s.p.x) == 'X')
    val ms = advance(xs, 'M')
    val as = advance(ms, 'A')
    val ss = advance(as, 'S')

    println(ss.size)
  }
  
}

object part2 {
  def run() = {
    val as = (
      for {
        x <- 0 to MaxX
        y <- 0 to MaxY if input(y)(x) == 'A' && x >= 1 && y >= 1 && x <= MaxX - 1 && y <= MaxY - 1
      } yield P(x, y)
    ).toList

    println(
      as.filter { p =>
        (
          (input(p.y - 1)(p.x - 1) == 'M' && input(p.y + 1)(p.x + 1) == 'S') ||
          (input(p.y - 1)(p.x - 1) == 'S' && input(p.y + 1)(p.x + 1) == 'M')
        ) &&
        (
          (input(p.y - 1)(p.x + 1) == 'M' && input(p.y + 1)(p.x - 1) == 'S') ||
          (input(p.y - 1)(p.x + 1) == 'S' && input(p.y + 1)(p.x - 1) == 'M')
        )
      }.size
    )
  }
}

part2.run()
