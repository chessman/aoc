case class Eq(testValue: Long, values: List[Long])
type Op = (v1: Long, v2: Long) => Long

val contestRaw = scala.io.Source
  .fromFile("input.txt")
  .getLines()
  .map { s =>
    val l = s.split(": ")
    Eq(l.head.toLong, l.tail.head.split(" ").map(_.toLong).toList)
  }
  .toList

val input = contestRaw

def check(ops: List[Op], testValue: Long, values: List[Long], res: Long): Boolean = {
  values match {
    case Nil => testValue == res
    case v :: l => ops.exists(op => check(ops, testValue, l, op(res, v)))
  }
}

def solve(ops: List[Op]) = {
  contestRaw.filter(eq => check(ops, eq.testValue, eq.values.tail, eq.values.head)).map(_.testValue).sum
}

object p1 {
  val ops: List[Op] = List(
    (v1: Long, v2: Long) => v1 + v2,
    (v1: Long, v2: Long) => v1 * v2
  )

  def run() = {
    println(solve(p1.ops))
  }
}

object p2 {
  val ops: List[Op] = List(
    (v1: Long, v2: Long) => v1 + v2,
    (v1: Long, v2: Long) => v1 * v2,
    (v1: Long, v2: Long) => s"$v1$v2".toLong
  )

  def run() = {
    println(solve(p2.ops))
  }
}

p2.run()
