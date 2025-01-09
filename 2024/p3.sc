object v1 {

  enum Symbol:
    case M, U, L, D, O, N, Apostrophe, T, ParL, ParR, Digit, Comma, Other

  object Symbol {
    def apply(c: Char) = c match {
      case 'm'  => Symbol.M
      case 'u'  => Symbol.U
      case 'l'  => Symbol.L
      case 'd'  => Symbol.D
      case 'o'  => Symbol.O
      case 'n'  => Symbol.N
      case '\'' => Symbol.Apostrophe
      case 't'  => Symbol.T
      case '('  => Symbol.ParL
      case '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' =>
        Symbol.Digit
      case ',' => Symbol.Comma
      case ')' => Symbol.ParR
      case _   => Symbol.Other
    }
  }

  case class State(
      firstNum: String,
      secondNum: String,
      first: Boolean,
      expect: List[Symbol],
      res: Int,
      enabled: Boolean,
      command: String
  ) {
    def setDigit(c: Char) = if (first) this.copy(firstNum = firstNum + c)
    else this.copy(secondNum = secondNum + c)
    def setCommand(c: Char) = this.copy(command = command + c)
    def next(s: Symbol) = this.copy(expect = expectNext(s, first))
    def done =
      command match {
        case "mul" =>
          if (firstNum.isEmpty() || secondNum.isEmpty() || !enabled)
            reinit
          else
            State.init(res + firstNum.toInt * secondNum.toInt, enabled)
        case "do" =>
          State.init(res, true)
        case "don't" =>
          State.init(res, false)
      }
    def firstNumDone = this.copy(first = false)
    def reinit = State.init(res, enabled)
  }

  object State {
    def init(res: Int, enabled: Boolean): State =
      State("", "", true, initialExpect, res, enabled, "")
    def init: State = init(0, true)
  }

  val initialExpect = List(Symbol.M, Symbol.D)

  def expectNext(got: Symbol, first: Boolean): List[Symbol] =
    got match
      case Symbol.M          => List(Symbol.U)
      case Symbol.U          => List(Symbol.L)
      case Symbol.L          => List(Symbol.ParL)
      case Symbol.D          => List(Symbol.O)
      case Symbol.O          => List(Symbol.N, Symbol.ParL)
      case Symbol.N          => List(Symbol.Apostrophe)
      case Symbol.Apostrophe => List(Symbol.T)
      case Symbol.T          => List(Symbol.ParL)
      case Symbol.ParL       => List(Symbol.Digit, Symbol.ParR)
      case Symbol.ParR       => List(Symbol.M)
      case Symbol.Digit =>
        if (first) List(Symbol.Digit, Symbol.Comma)
        else List(Symbol.Digit, Symbol.ParR)
      case Symbol.Comma => List(Symbol.Digit)
      case Symbol.Other => initialExpect

  def enter(c: Char, s: Symbol, state: State): State = {
    (s match {
      case Symbol.M | Symbol.U | Symbol.L | Symbol.D | Symbol.O | Symbol.N |
          Symbol.Apostrophe | Symbol.T =>
        state.setCommand(c)
      case Symbol.Digit => state.setDigit(c)
      case Symbol.Comma => state.firstNumDone
      case Symbol.ParR  => state.done
      case _            => state
    }).next(s)
  }

  def parse(c: Char, state: State): State = {
    val s = Symbol(c)
    if (state.expect.contains(s)) {
      enter(c, s, state)
    } else if (initialExpect.contains(s)) {
      enter(c, s, state.reinit)
    } else state.reinit
  }

  def parseAll(input: String, state: State = State.init): State = {
    if (input.size == 0) state
    else parseAll(input.drop(1), parse(input.head, state))
  }
}

val input = scala.io.Source
  .fromFile("input.txt")
  .mkString

println(v1.parseAll(input))
