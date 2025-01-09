val contestRaw = scala.io.Source
  .fromFile("input.txt")
  .getLines()
  .next()

val demoRaw = "2333133121414131402"

val input = contestRaw
var lastFileIdx = (input.size - 1) - (input.size - 1) % 2

def subsum(index: Long, fileno: Long, count: Long): Long =
  fileno * count * (2 * index + (count - 1)) / 2

object p1 {
  def compactAndChecksum(s: String) = {
    var p1 = 0
    var p2 = lastFileIdx
    var checksum = 0L
    var index = 0L
    var leftoverFile = 0

    while (p1 <= p2 || leftoverFile > 0) {
      if (p1 % 2 == 0 && p1 <= p2) {
        val count = input(p1).asDigit
        val fileno = p1 / 2
        checksum += subsum(index, fileno, count)
        // println(s"A p1: $p1, p2: $p2, index: $index, fileno: $fileno, count: $count, sum: ${subsum(index, fileno, count)}")
        index += count
      } else {
        var free = input(p1).asDigit
        while (free > 0 && (p2 >= p1 || leftoverFile > 0)) {
          if (leftoverFile == 0) {
            leftoverFile = input(p2).asDigit
            p2 = p2 - 2
          }
          val fileno = p2 / 2 + 1
          val count = Math.min(leftoverFile, free)
          checksum += subsum(index, fileno, count)
          // println(s"B p1: $p1, p2: $p2, index: $index, fileno: $fileno, count: $count, sum: ${subsum(index, fileno, count)}")
          index += count
          free -= count
          leftoverFile -= count
        }
      }
      p1 += 1
    }
    checksum
  }
  def run() = {
    println(compactAndChecksum(input))
  }
}

object p2 {
  enum Fs:
    case File(fileno: Long, size: Long)
    case Space(size: Long, files: Vector[File] = Vector.empty)

  def parse(input: String) = for {
    i <- (0 until input.size).toVector
  } yield
    if (i % 2 == 0) Fs.File(i / 2, input(i).asDigit)
    else Fs.Space(input(i).asDigit)

  def defrag(s: String) = {
    (lastFileIdx to 0 by -2).foldLeft(parse(input)) { case (files, i) =>
      val file = files(i).asInstanceOf[Fs.File]
      val freeSpaceIdx =
        (1 to i by 2).find(i =>
          files(i).asInstanceOf[Fs.Space].size >= file.size
        )
      freeSpaceIdx match {
        case None => files
        case Some(idx) =>
          val freeSpaceToMerge = files(i - 1).asInstanceOf[Fs.Space]
          val freeSpace = files(idx).asInstanceOf[Fs.Space]
          files
            .patch(
              i - 1,
              Vector(
                freeSpaceToMerge.copy(size = freeSpaceToMerge.size + file.size)
              ),
              2
            )
            .updated(
              idx,
              freeSpace.copy(
                size = freeSpace.size - file.size,
                files = freeSpace.files.appended(file)
              )
            )
      }
    }
  }

  case class Res(sum: Long, index: Long)

  def checksum(files: Vector[Fs], res: Res = Res(0, 0)): Res = {
    files.foldLeft(res) { case (res, file) =>
      file match {
        case Fs.File(fileno, size) =>
          Res(res.sum + subsum(res.index, fileno, size), res.index + size)
        case Fs.Space(size, files) =>
          val newRes = checksum(files, res)
          newRes.copy(index = newRes.index + size)
      }
    }
  }

  def run() = {
    println(checksum(defrag(input)))
  }
}

p2.run()
