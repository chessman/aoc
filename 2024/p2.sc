object v1 {
  def valid(x: Int, gt: Boolean) =
    x.abs >= 1 && x.abs <= 3 && (if (gt) x > 0 else x < 0)

  def diff(input: List[Int]) = input.dropRight(1).zip(input.drop(1)).map(_ - _)

  def isValid(x: List[Int]) = x.forall(a => valid(a, x.head > 0))

  def genOpts(input: List[Int]): List[List[Int]] =
    input :: (0 until input.size).map(s => input.patch(s, Nil, 1)).toList

  def safe(input: List[Int]) =
    genOpts(input).exists(x => isValid(diff(x)))
}

val input = scala.io.Source
  .fromFile("input.txt")
  .mkString
  .split("\n")
  .map(_.split(" ").map(_.toInt).toList)
  .toList

val v1Solutions = input.map(v1.safe).filter(identity)

println(v1Solutions.size)
