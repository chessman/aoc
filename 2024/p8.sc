val contestRaw = scala.io.Source
  .fromFile("input.txt")
  .getLines()
  .toArray

val demoRaw = """
............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............
""".split("\n").drop(1)

val input = contestRaw

case class Coord(down: Int, right: Int) {
  def antinodes(c: Coord) = {
    // distances:
    // 2 * down - c.down - down = down - c.down
    // 2 * down - c.down - c.down = 2 * (down - c.down)
    List(
      Coord(2 * down - c.down, 2 * right - c.right),
      Coord(2 * c.down - down, 2 * c.right - right)
    ).filter(c => c.down >= 0 && c.down < input.size && c.right >= 0 && c.right < input(0).size)
  }

  def distantAntinodes(c: Coord) = {
    (1 to input.size).flatMap { m =>
      List(
        Coord(m * down - (m - 1) * c.down, m * right - (m - 1) * c.right),
        Coord(m * c.down - (m - 1) * down, m * c.right - (m - 1) * right)
      ).filter(c => c.down >= 0 && c.down < input.size && c.right >= 0 && c.right < input(0).size)
    }
  }
}

val antennas: List[(Char, Seq[Coord])] =
  (
    for {
      i <- 0 until input.size
      j <- 0 until input(0).size if input(i)(j) != '.'
    } yield input(i)(j) -> Coord(i, j)
  ).groupMap(_._1)(_._2).toList

def genPairs(antennas: List[Coord]): List[(Coord, Coord)] = {
  antennas match {
    case Nil => Nil
    case a :: l => l.map { b => (a, b) } ::: genPairs(l)
  }
}

object p1 {
  def countAntinodes(antennas: List[(Char, Seq[Coord])]): Int = {
    val pairs = antennas.flatMap(a => genPairs(a._2.toList))
    pairs.flatMap { case (c1, c2) => c1.antinodes(c2) }.distinct.size
  }

  def run() = {
    println(countAntinodes(antennas))
  }
}

object p2 {
  def countDistantAntinodes(antennas: List[(Char, Seq[Coord])]): Int = {
    val pairs = antennas.flatMap(a => genPairs(a._2.toList))
    pairs.flatMap { case (c1, c2) => c1.distantAntinodes(c2) }.distinct.size
  }

  def run() = {
    println(countDistantAntinodes(antennas))
  }
}

p2.run()

