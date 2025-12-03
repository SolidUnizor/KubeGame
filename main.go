package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GridSize                = 50
	GridCols                = 12
	GridRows                = 8
	ScreenWidth             = GridCols * GridSize
	ScreenHeight            = GridRows * GridSize
	FLAG_WINDOW_RESIZABLE   = 0x00000004
	FLAG_WINDOW_UNDECORATED = 0x00000008
)

// Die представляет кубик с отслеживанием всех сторон
type Die struct {
	Top, Bottom, Front, Back, Left, Right int // все 6 сторон
	CurrentTop                            int // число на верхней стороне
}

// Direction направление движения
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Player представляет игрока-кубик
type Player struct {
	X, Y int // позиция в сетке
	Die  Die // кубик с отслеживанием сторон
}

// Level представляет уровень игры
type Level struct {
	Player Player
	Finish struct {
		X, Y   int
		Number int // число на финише
	}
	Walls []struct {
		X, Y int
	}
	Won bool
}

// NewDie создает новый кубик
func NewDie() Die {
	// Начинаем с кубика, где верх = 1
	die := Die{
		Top:    1,
		Bottom: 6, // противоположная сторона
		Front:  2,
		Back:   5, // противоположная сторона
		Left:   3,
		Right:  4, // противоположная сторона
	}
	die.CurrentTop = die.Top
	return die
}

// Roll перекатывает кубик в указанном направлении
func (d *Die) Roll(dir Direction) {
	switch dir {
	case Up:
		// При движении вверх кубик катится назад
		newTop := d.Front
		d.Front = d.Bottom
		d.Bottom = d.Back
		d.Back = d.Top
		d.Top = newTop
	case Down:
		// При движении вниз кубик катится вперед
		newTop := d.Back
		d.Back = d.Bottom
		d.Bottom = d.Front
		d.Front = d.Top
		d.Top = newTop
	case Left:
		// При движении влево кубик катится направо
		newTop := d.Right
		d.Right = d.Bottom
		d.Bottom = d.Left
		d.Left = d.Top
		d.Top = newTop
	case Right:
		// При движении вправо кубик катится налево
		newTop := d.Left
		d.Left = d.Bottom
		d.Bottom = d.Right
		d.Right = d.Top
		d.Top = newTop
	}
	d.CurrentTop = d.Top
}

// NewPlayer создает нового игрока
func NewPlayer(x, y int) Player {
	return Player{
		X:   x,
		Y:   y,
		Die: NewDie(),
	}
}

// Move двигает игрока в указанном направлении
func (p *Player) Move(dx, dy int, dir Direction) {
	p.X += dx
	p.Y += dy
	p.Die.Roll(dir)
}

// NewLevel создает новый уровень
func NewLevel() Level {
	level := Level{}

	// Создаем игрока
	level.Player = NewPlayer(1, 1)

	// Создаем финишную клетку
	level.Finish.X = GridCols - 2
	level.Finish.Y = GridRows - 2
	level.Finish.Number = rand.Intn(6) + 1

	// Создаем стены (простая лабиринтоподобная структура)
	level.Walls = []struct {
		X, Y int
	}{
		// Внешние стены
		{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, {6, 0}, {7, 0}, {8, 0}, {9, 0}, {10, 0}, {11, 0},
		{0, 7}, {1, 7}, {2, 7}, {3, 7}, {4, 7}, {5, 7}, {6, 7}, {7, 7}, {8, 7}, {9, 7}, {10, 7}, {11, 7},
		{0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6},
		{11, 1}, {11, 2}, {11, 3}, {11, 4}, {11, 5}, {11, 6},

		// Внутренние стены
		{3, 2}, {3, 3}, {3, 4},
		{7, 4}, {7, 5}, {7, 6},
		{5, 1}, {6, 1}, {8, 1},
		{2, 5}, {4, 5}, {9, 5},
	}

	return level
}

// IsWall проверяет, является ли клетка стеной
func (l *Level) IsWall(x, y int) bool {
	for _, wall := range l.Walls {
		if wall.X == x && wall.Y == y {
			return true
		}
	}
	return false
}

// IsValidMove проверяет, можно ли двигаться в указанную клетку
func (l *Level) IsValidMove(x, y int) bool {
	// Проверяем границы сетки
	if x < 0 || x >= GridCols || y < 0 || y >= GridRows {
		return false
	}

	// Проверяем стены
	if l.IsWall(x, y) {
		return false
	}

	return true
}

// CheckWin проверяет условие победы
func (l *Level) CheckWin() bool {
	return l.Player.X == l.Finish.X &&
		l.Player.Y == l.Finish.Y &&
		l.Player.Die.CurrentTop == l.Finish.Number
}

// DrawDie рисует кубик с числом
func DrawDie(x, y int, number int) {
	// Размер кубика
	size := 40

	// Позиция для рисования
	drawX := x*GridSize + (GridSize-size)/2
	drawY := y*GridSize + (GridSize-size)/2

	// Цвет кубика зависит от числа
	var color rl.Color
	switch number {
	case 1:
		color = rl.Red
	case 2:
		color = rl.Orange
	case 3:
		color = rl.Yellow
	case 4:
		color = rl.Green
	case 5:
		color = rl.Blue
	case 6:
		color = rl.Purple
	default:
		color = rl.Gray
	}

	// Рисуем кубик
	rl.DrawRectangle(int32(drawX), int32(drawY), int32(size), int32(size), color)
	rl.DrawRectangleLines(int32(drawX), int32(drawY), int32(size), int32(size), rl.Black)

	// Рисуем число на кубике
	text := string(rune('0' + number))
	fontSize := int32(24)
	textX := drawX + (size-int(rl.MeasureText(text, fontSize)))/2
	textY := drawY + (size-24)/2

	rl.DrawText(text, int32(textX), int32(textY), fontSize, rl.White)
}

// DrawGrid рисует игровую сетку
func DrawGrid() {
	// Рисуем фон сетки
	for y := 0; y < GridRows; y++ {
		for x := range GridCols {
			// Чередующиеся цвета клеток
			var color rl.Color
			if (x+y)%2 == 0 {
				color = rl.LightGray
			} else {
				color = rl.Gray
			}

			rl.DrawRectangle(int32(x*GridSize), int32(y*GridSize), int32(GridSize), int32(GridSize), color)
			rl.DrawRectangleLines(int32(x*GridSize), int32(y*GridSize), int32(GridSize), int32(GridSize), rl.Black)
		}
	}
}

// DrawWalls рисует стены
func DrawWalls(walls []struct{ X, Y int }) {
	for _, wall := range walls {
		x := wall.X * GridSize
		y := wall.Y * GridSize
		rl.DrawRectangle(int32(x), int32(y), int32(GridSize), int32(GridSize), rl.Brown)
		rl.DrawRectangleLines(int32(x), int32(y), int32(GridSize), int32(GridSize), rl.Black)
	}
}

// DrawFinish рисует финишную клетку
func DrawFinish(finish struct{ X, Y, Number int }) {
	x := finish.X * GridSize
	y := finish.Y * GridSize

	// Рисуем финишную клетку
	rl.DrawRectangle(int32(x), int32(y), int32(GridSize), int32(GridSize), rl.Gold)
	rl.DrawRectangleLines(int32(x), int32(y), int32(GridSize), int32(GridSize), rl.Black)

	// Рисуем число на финише
	text := string(rune('0' + finish.Number))
	fontSize := int32(24)
	textWidth := rl.MeasureText(text, fontSize)
	textX := x + (GridSize-int(textWidth))/2
	textY := y + (GridSize-24)/2

	rl.DrawText(text, int32(textX), int32(textY), fontSize, rl.Black)
}

// DrawUI рисует пользовательский интерфейс
func DrawUI(player Player, finishNumber int, won bool) {
	// Информация о текущем числе кубика
	text := string(rune('0' + player.Die.CurrentTop))
	rl.DrawText(text, 10, 10, 20, rl.Black)

	// Информация о числе на финише
	finishText := string(rune('0' + finishNumber))
	rl.DrawText(finishText, 10, 40, 20, rl.Black)

	// Инструкции
	instructions := "Use WASD or Arrow Keys to move"
	rl.DrawText(instructions, 10, 70, 20, rl.Black)

	// Сообщение о победе
	if won {
		winText := "YOU WIN! Press R to restart"
		textWidth := rl.MeasureText(winText, 30)
		textX := (ScreenWidth - int(textWidth)) / 2
		textY := ScreenHeight/2 - 15

		rl.DrawRectangle(int32(textX-10), int32(textY-10), int32(textWidth+20), 60, rl.Green)
		rl.DrawText(winText, int32(textX), int32(textY), 30, rl.White)
	}
}

// HandleInput обрабатывает ввод игрока
func HandleInput(level *Level) {
	if level.Won {
		// Если игра выиграна, только перезапуск
		if rl.IsKeyPressed(rl.KeyR) {
			*level = NewLevel()
		}
		return
	}

	// Движение по WASD
	if rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp) {
		newX, newY := level.Player.X, level.Player.Y-1
		if level.IsValidMove(newX, newY) {
			level.Player.Move(0, -1, Up)
		}
	}
	if rl.IsKeyPressed(rl.KeyS) || rl.IsKeyPressed(rl.KeyDown) {
		newX, newY := level.Player.X, level.Player.Y+1
		if level.IsValidMove(newX, newY) {
			level.Player.Move(0, 1, Down)
		}
	}
	if rl.IsKeyPressed(rl.KeyA) || rl.IsKeyPressed(rl.KeyLeft) {
		newX, newY := level.Player.X-1, level.Player.Y
		if level.IsValidMove(newX, newY) {
			level.Player.Move(-1, 0, Left)
		}
	}
	if rl.IsKeyPressed(rl.KeyD) || rl.IsKeyPressed(rl.KeyRight) {
		newX, newY := level.Player.X+1, level.Player.Y
		if level.IsValidMove(newX, newY) {
			level.Player.Move(1, 0, Right)
		}
	}
}

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Создаем окно
	rl.InitWindow(ScreenWidth, ScreenHeight, "KubeGame - Die Rolling Puzzle")
	rl.SetWindowState(FLAG_WINDOW_RESIZABLE)
	rl.SetWindowState(FLAG_WINDOW_UNDECORATED)
	rl.SetTargetFPS(60)

	// Создаем уровень
	level := NewLevel()

	// Главный игровой цикл
	for !rl.WindowShouldClose() {
		// Обновление
		HandleInput(&level)

		// Проверяем победу
		if !level.Won && level.CheckWin() {
			level.Won = true
		}

		// Рендеринг
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Рисуем игровое поле
		DrawGrid()
		DrawWalls(level.Walls)
		DrawFinish(level.Finish)
		DrawDie(level.Player.X, level.Player.Y, level.Player.Die.CurrentTop)

		// Рисуем UI
		DrawUI(level.Player, level.Finish.Number, level.Won)

		rl.EndDrawing()
	}

	// Закрываем окно
	rl.CloseWindow()
}
