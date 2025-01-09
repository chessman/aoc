val contestRaw = scala.io.Source
  .fromFile("input.txt")
  .getLines()
  .map(_.split("").map(_.head))
  .toArray

val demoRaw = """
....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...""".split("\n").drop(1).map(_.split("").map(_.head)).toArray

val input = contestRaw

val iSize = input.size
val jSize = input(0).size

val start = (0 until iSize)
  .flatMap(i => (0 until jSize).find(j => input(i)(j) == '^').map(j => (i, j)))
  .head
//input(45).update(47, '.')

type Dir = Int
case class P(pos: (Int, Int), dir: Dir)

val dirs = List(
      (p: Tuple2[Int, Int]) => (p._1 - 1, p._2),
      (p: Tuple2[Int, Int]) => (p._1, p._2 + 1),
      (p: Tuple2[Int, Int]) => (p._1 + 1, p._2),
      (p: Tuple2[Int, Int]) => (p._1, p._2 - 1)
    )

def loop(
  positions: Vector[P]
): Vector[P] = {
  val dir = dirs(positions.head.dir)
  val newPos = dir(positions.head.pos)
  if (
    newPos._1 < 0 || newPos._1 >= iSize || newPos._2 < 0 || newPos._2 >= jSize || positions.size >= iSize * jSize
  )
    positions
  else {
    if (input(newPos._1)(newPos._2) == '#') loop(P(positions.head.pos, (positions.head.dir + 1) % 4) +: positions.drop(1))
    else {
      loop(P(newPos, positions.head.dir) +: positions)
    }
  }
}

object p1 {
  def run() = {
    println(loop(Vector(P(start, 0))).map(_.pos).distinct.size)
  }
}

object p2 {
  def run() = {
    val positions = loop(Vector(P(start, 0)))
    val res = positions.dropRight(1).distinctBy(_.pos).count { obstacle =>
      input(obstacle.pos._1).update(obstacle.pos._2, '#')
      val res = loop(Vector(P(start, 0))).size
      input(obstacle.pos._1).update(obstacle.pos._2, '.')
      res == iSize * jSize
    }

    println(res)
  }
}

p2.run()
